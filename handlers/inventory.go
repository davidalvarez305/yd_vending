package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/types"
)

func createInventoryContext() map[string]any {
	return map[string]any{
		"PageTitle":       constants.CompanyName,
		"MetaDescription": "Get a quote for vending machine services.",
		"SiteName":        constants.SiteName,
		"StaticPath":      constants.StaticPath,
		"MediaPath":       constants.MediaPath,
		"PhoneNumber":     constants.DavidPhoneNumber,
		"CurrentYear":     time.Now().Year(),
		"CompanyName":     constants.CompanyName,
	}
}

func InventoryHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx := createInventoryContext()
	ctx["PagePath"] = constants.RootDomain + path

	stats, err := database.GetDashboardStats()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting dashboard stats from DB.", http.StatusInternalServerError)
		return
	}
	ctx["DashboardStats"] = stats

	switch r.Method {
	case http.MethodPost:
		switch path {
		case "/inventory/product":
			PostProduct(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodGet:
		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				GetProductDetail(w, r, ctx)
				return
			}
		}

		switch path {
		case "/inventory/product":
			GetProducts(w, r, ctx)
		case "/inventory/transaction":
			GetTransactions(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				PutProduct(w, r)
				return
			}
		}

		switch path {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				DeleteProduct(w, r)
				return
			}
		}

		switch path {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.INVENTORY_TEMPLATES_DIR + "products.html"
	createProductForm := constants.INVENTORY_TEMPLATES_DIR + "create_product_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "products_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createProductForm}

	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	categories, err := database.GetProductCategories()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting product categories from DB.", http.StatusInternalServerError)
		return
	}

	products, totalRows, err := database.GetProductList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting products from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Products — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Products"] = products
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum
	data["ProductCategories"] = categories

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.ProductForm
	err = decoder.Decode(&form, r.PostForm)

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error decoding form data.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.CreateProduct(form)
	if err != nil {
		fmt.Printf("Error creating product: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create product.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1 // Always default to one after new entity is created
	products, totalRows, err := database.GetProductList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting products from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "products_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "products_table.html",
		Data: map[string]any{
			"Products":    products,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.ProductForm
	err = decoder.Decode(&form, r.PostForm)

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error decoding form data.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productId, err := helpers.GetFirstIDAfterPrefix(r, "/inventory/product/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateProduct(productId, form)
	if err != nil {
		fmt.Printf("Error updating product: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update product.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Product updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetProductDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "product_detail.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.INVENTORY_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	productId := strings.TrimPrefix(r.URL.Path, "/inventory/product/")

	productDetails, err := database.GetProductDetails(productId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting product details from DB.", http.StatusInternalServerError)
		return
	}

	productCategories, err := database.GetProductCategories()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting product categories.", http.StatusInternalServerError)
		return
	}

	suppliers, err := database.GetSuppliers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting suppliers.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Product Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Product"] = productDetails
	data["ProductCategories"] = productCategories
	data["Suppliers"] = suppliers

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := helpers.GetFirstIDAfterPrefix(r, "/inventory/product/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteProduct(productId)
	if err != nil {
		fmt.Printf("Error deleting product: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete product.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	products, totalRows, err := database.GetProductList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting products from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "products_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "products_table.html",
		Data: map[string]any{
			"Products":    products,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetTransactions(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.INVENTORY_TEMPLATES_DIR + "transactions.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "transactions_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table}

	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	var params types.GetTransactionsParams
	params.TransactionType = helpers.SafeStringToPointer(r.URL.Query().Get("transaction_type"))
	params.Machine = helpers.SafeStringToPointer(r.URL.Query().Get("machine"))
	params.Location = helpers.SafeStringToPointer(r.URL.Query().Get("location"))
	params.Product = helpers.SafeStringToPointer(r.URL.Query().Get("product"))
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))
	params.DateFrom = helpers.SafeStringToInt64Pointer(r.URL.Query().Get("date_from"))
	params.DateTo = helpers.SafeStringToInt64Pointer(r.URL.Query().Get("date_to"))

	transactions, totalRows, err := database.GetTransactionList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting transactions from DB.", http.StatusInternalServerError)
		return
	}

	machines, err := database.GetMachines()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machines from DB.", http.StatusInternalServerError)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting locations from DB.", http.StatusInternalServerError)
		return
	}

	products, err := database.GetProducts()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting products from DB.", http.StatusInternalServerError)
		return
	}

	transactionTypes, err := database.GetTransactionTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting transaction types from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Transactions — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Transactions"] = transactions
	data["Machines"] = machines
	data["Locations"] = locations
	data["Products"] = products
	data["TransactionTypes"] = transactionTypes
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
