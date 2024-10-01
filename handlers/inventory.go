package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
)

var inventoryBaseFilePath = constants.INVENTORY_TEMPLATES_DIR + "base.html"
var inventoryFooterFilePath = constants.INVENTORY_TEMPLATES_DIR + "footer.html"

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
	ctx := createCrmContext()
	ctx["PagePath"] = constants.RootDomain + r.URL.Path

	stats, err := database.GetDashboardStats()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting dashboard stats from DB.", http.StatusInternalServerError)
		return
	}
	ctx["DashboardStats"] = stats

	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/inventory/product":
			PostProduct(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodGet:
		switch r.URL.Path {
		case "/inventory/product":
			GetProducts(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		switch r.URL.Path {
		case "/inventory/product":
			PutProduct(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodDelete:
		switch r.URL.Path {
		case "/inventory/product":
			DeleteProduct(w, r)
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
