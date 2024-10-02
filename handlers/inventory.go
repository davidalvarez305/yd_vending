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
		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				PostProductBatch(w, r)
				return
			}
			return
		}

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
			return
		}

		switch path {
		case "/inventory/product":
			GetProducts(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				PutProductBatch(w, r)
				return
			}
			return
		}

		switch path {
		case "/inventory/product":
			PutProduct(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/inventory/product/") {
			parts := strings.Split(path, "/")
			if len(parts) >= 6 && parts[4] == "batch" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				DeleteProductBatch(w, r)
				return
			}
			return
		}

		if strings.HasPrefix(path, "/inventory/product/") {
			if len(path) > len("/inventory/product/") && helpers.IsNumeric(path[len("/inventory/product/"):]) {
				DeleteProduct(w, r)
				return
			}
			return
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
	hasPageNum := r.URL.Query().Has("pageNum")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("pageNum"))
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

func DeleteProductBatch(w http.ResponseWriter, r *http.Request) {
	productBatchId, err := helpers.GetSecondIDFromPath(r, "/inventory/product/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteProductBatch(productBatchId)
	if err != nil {
		fmt.Printf("Error deleting product batch: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete product batch.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productBatches, err := database.GetProductBatchList(string(productBatchId))
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
			"ProductBatches": productBatches,
			"CurrentPage":    1,
			"MaxPages":       helpers.CalculateMaxPages(len(productBatches), constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetProductDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "product_detail.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
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

	productBatches, err := database.GetProductBatchList(productId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting product batches.", http.StatusInternalServerError)
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
	data["ProductBatches"] = productBatches
	data["Suppliers"] = suppliers

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostProductBatch(w http.ResponseWriter, r *http.Request) {
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

	var form types.ProductBatchForm
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

	err = database.CreateProductBatch(form)
	if err != nil {
		fmt.Printf("Error creating product batch: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create product batch.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productId, err := helpers.GetFirstIDAfterPrefix(r, "/inventory/product/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid URL.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productBatches, err := database.GetProductBatchList(string(productId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting product batches from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "product_batches_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "product_batches_table.html",
		Data: map[string]any{
			"ProductBatches": productBatches,
			"CurrentPage":    1,
			"MaxPages":       helpers.CalculateMaxPages(len(productBatches), constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutProductBatch(w http.ResponseWriter, r *http.Request) {
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

	var form types.ProductBatchForm
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
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid URL.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.UpdateProductBatch(productId, form)
	if err != nil {
		fmt.Printf("Error updating product batch: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update product batch.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productBatches, err := database.GetProductBatchList(string(productId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting product batches from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "product_batches_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "product_batches_table.html",
		Data: map[string]any{
			"ProductBatches": productBatches,
			"CurrentPage":    1,
			"MaxPages":       helpers.CalculateMaxPages(len(productBatches), constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
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
				"Message": "Error getting product batches from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "product_batches_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "product_batches_table.html",
		Data: map[string]any{
			"ProductBatches": products,
			"CurrentPage":    pageNum,
			"MaxPages":       helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
