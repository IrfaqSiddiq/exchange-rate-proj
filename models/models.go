package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"project_first/config"
	"project_first/utility"
	"time"
)

type ExchangeDataOpenExchangeRate struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

type ItemsDetailList struct {
	Item            string  `json:"item"`
	USDCost         float64 `json:"usd_cost"`
	ActualZMWCost   float64 `json:"actual_zmw_cost"`
	LatestZMWCost   float64 `json:"latest_zmw_cost"`
	ExchangeAmount  float64 `json:"exchange_rate"`
	PurchaseDate    string  `json:"purchase_date"`
	USDSellingPrice float64 `json:"usd_selling_price"`
	ZMWSellingPrice float64 `json:"zmw_selling_price"`
	ProfitPerc      float64 `json:"profit_perc"`
	LossPrevented   float64 `json:"loss_prevented"`
}

type AdminSetting struct {
	ProfitPerc float64 `json:"profit_perc"`
}

func InsertItems(item string, amount float64, purchasingDate string) error {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("InsertItems: Failed while connecting with the database :", err)
		return err
	}
	defer db.Close()
	query := `INSERT INTO items_info(item,amount,purchase_date)VALUES($1,$2,$3)`
	_, err = db.Exec(query, item, amount, purchasingDate)
	if err != nil {
		log.Println("InsertItems: Failed while executing the query with error :", err)
		return err
	}
	return nil
}

func DisplayAllItems() ([]ItemsDetailList, error) {
	var itemsList []ItemsDetailList
	db, err := config.GetDB2()
	if err != nil {
		log.Println("DisplayAllItems: Failed while connecting with the database :", err)
		return itemsList, err
	}
	defer db.Close()
	adminSetting, err := GetAdminSettingValues()
	if err != nil {
		log.Println("DisplayAllItems: failed while fetching profit percentage with error:", err)
		return itemsList, err
	}
	query := `SELECT 
				ii.item,
				ii.amount,
				ii.purchase_date,
				er.exchange_amount,
				(SELECT exchange_amount FROM exchange_rates WHERE DATE(exchange_rate_time) = '` + "2024-02-01" + `')
			FROM
				items_info as ii, exchange_rates as er
			WHERE
				DATE(ii.purchase_date)=DATE(er.exchange_rate_time)`
	fmt.Println("**query", query)
	rows, err := db.Query(query)
	if err != nil {
		log.Println("DisplayAllItems: Failed while executing the query with error :", err)
		return itemsList, err
	}
	for rows.Next() {
		var (
			item                 sql.NullString
			amount               sql.NullFloat64
			purchaseDate         sql.NullTime
			exchangeAmount       sql.NullFloat64
			latestExchangeAmount sql.NullFloat64
		)
		err = rows.Scan(&item, &amount, &purchaseDate, &exchangeAmount, &latestExchangeAmount)
		if err != nil {
			log.Println("DisplayAllItems: Failed while scanning the query with error :", err)
			continue
		}
		costInZMW := amount.Float64 * exchangeAmount.Float64
		latestZmwCost := amount.Float64 * latestExchangeAmount.Float64
		sellingPriceInUSD := amount.Float64 + (amount.Float64 * adminSetting.ProfitPerc / 100)
		sellingPriceInZMW := latestZmwCost + (latestZmwCost * adminSetting.ProfitPerc / 100)
		defaultSellingPrice := costInZMW + (costInZMW * adminSetting.ProfitPerc / 100)
		lossPrevented := sellingPriceInZMW - defaultSellingPrice

		itemsList = append(itemsList, ItemsDetailList{
			Item:            item.String,
			USDCost:         amount.Float64,
			PurchaseDate:    purchaseDate.Time.Format("2006-01-02"),
			ExchangeAmount:  latestExchangeAmount.Float64,
			ActualZMWCost:   costInZMW,
			LatestZMWCost:   latestZmwCost,
			USDSellingPrice: sellingPriceInUSD,
			ZMWSellingPrice: math.Floor(sellingPriceInZMW*100) / 100,
			ProfitPerc:      adminSetting.ProfitPerc,
			LossPrevented:   lossPrevented,
		})
	}
	return itemsList, nil
}

// input: exchangeCurrency string
// output: (float64, int64, error)
// func GetExchangeAmountOpenExchange retrieves the exchange amount for a specific currency
// from an API and returns it as a float64 value along with any error encountered.
func GetExchangeAmountOpenExchange(exchangeCurrency string) (float64, int64, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	accessKey := utility.KeyForCurrencyExchangeOpenExchange()
	reqUrl := "https://openexchangerates.org/api/latest.json"
	URL, err := url.Parse(reqUrl)
	if err != nil {
		log.Println("GetExchangeAmountOpenExchange: failed to parse with error : ", err)
		return 0.0, 0, err
	}
	parameters := url.Values{}
	parameters.Add("app_id", accessKey)
	parameters.Add("symbols", exchangeCurrency)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	response, err := client.Get(url)
	if err != nil {
		log.Println("GetExchangeAmountOpenExchange: failed to get data from openexchangerates.org ", err)
		return 0.0, 0, err
	}
	fmt.Println("**url************", url)
	exchangeData := ExchangeDataOpenExchangeRate{}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("GetExchangeAmountOpenExchange: failed to get response data with error: ", err)
		return 0.0, 0, err
	}
	// fmt.Println("*********data", data)
	err = json.Unmarshal(data, &exchangeData)
	if err != nil {
		fmt.Println("GetExchangeAmountOpenExchange: failed: while unmarshal with error: ", err)
		// err = errors.New(string(data))
		return 0.0, 0, err
	}
	fmt.Println("****amount in models", exchangeData.Rates["ZMW"])
	return exchangeData.Rates["ZMW"], exchangeData.Timestamp, nil
}

// input:exchangeCurrency
// output:int64, error
// GetSupportedCountryID will fetch country id corresponding to exchangeCurrency
func GetSupportedCountryID(exchangeCurrency string) (int64, error) {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("GetSupportedCountryID: failed connecting to the database with error:", err)
		return 0, err
	}
	defer db.Close()
	var id sql.NullInt64
	query := `SELECT id FROM supported_countries WHERE currency_code = $1`
	err = db.QueryRow(query, exchangeCurrency).Scan(&id)
	if err != nil {
		log.Println("GetSupportedCountryID: failed while executing the query:", err)
		return 0, err
	}
	return id.Int64, nil
}

// input: exchangeAmount
// output: error
// func StoreExchangedAmount stores the provided exchange amount in a database table.
func StoreExchangedAmount(exchangeAmount float64, timeInUnix, countryId int64) error {
	var err error
	db, err := config.GetDB2()
	if err != nil {
		log.Println("StoreExchangedAmount: failed connecting to the database with error:", err)
		return err
	}
	defer db.Close()
	var (
		utcTime       time.Time
		utcTimeString string
	)
	if timeInUnix != 0 {
		utcTime = time.Unix(timeInUnix, 0).UTC()
		utcTimeString = utcTime.Format("2006-01-02")
	} else {
		utcTimeString = time.Now().Format("2006-01-02")
	}
	query := `INSERT INTO exchange_rates(exchange_amount,exchange_rate_time,created_at,country_id)VALUES($1,$2,NOW(),$3)`
	_, err = db.Exec(query, exchangeAmount, utcTimeString, countryId)
	if err != nil {
		log.Println("StoreExchangedAmount: failed while executing the query with error:", err)
		return err
	}
	return nil
}

func GetAdminSettingValues() (AdminSetting, error) {
	var adminSetting AdminSetting
	db, err := config.GetDB2()
	if err != nil {
		log.Println("GetAdminSettingValues: failed connecting to the database with error:", err)
		return adminSetting, err
	}
	defer db.Close()
	var profitPercentage sql.NullFloat64
	query := `SELECT profit_perc FROM admin_settings`
	err = db.QueryRow(query).Scan(&profitPercentage)
	if err != nil {
		log.Println("GetAdminSettingValues: failed while executing the query with error:", err)
		return adminSetting, err
	}
	adminSetting = AdminSetting{
		ProfitPerc: profitPercentage.Float64,
	}
	return adminSetting, nil
}

func UpdateAdminSettings(profitPerc float64) error {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("UpdateAdminSettings: failed connecting to the database with error:", err)
		return err
	}
	defer db.Close()
	query := `UPDATE admin_settings SET profit_perc = $1 `
	_, err = db.Exec(query, profitPerc)
	if err != nil {
		log.Println("UpdateAdminSettings: failed while executing the query with error:", err)
		return err
	}
	return nil
}
