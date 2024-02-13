package controllers

import (
	"fmt"
	"log"
	"net/http"
	"project_first/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DisplayInsertionModule(c *gin.Context) {
	c.HTML(http.StatusOK, "insert_module.html", gin.H{})
}

func DisplayItems(c *gin.Context) {
	c.HTML(http.StatusOK, "display_items.html", gin.H{})
}

func GetAdminSettings(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_settings.html", gin.H{})
}

func InsertItems(c *gin.Context) {
	item := c.PostForm("item")
	price := c.PostForm("amount")
	fmt.Println("****price", price)
	amount, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("InsertItems: Error while converting string to float:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Error while converting string to float",
			"error":   err,
		})
		return
	}
	purchasingDate := c.PostForm("date")
	fmt.Println("*****purchasing_date", purchasingDate)
	layout := "2006-01-02"
	_, err = time.Parse(layout, purchasingDate)
	if err != nil {
		log.Println("InsertItems: purchasing date must be of date datatype: error is", err)
		purchasingDate = time.Now().Format("2006-01-02T15:04:05Z")
	}
	err = models.InsertItems(item, amount, purchasingDate)
	if err != nil {
		fmt.Println("InsertItems: failed while inserting records in database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed while inserting records ",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "successfully inserted records",
	})
}

func DisplayAllItems(c *gin.Context) {
	itemsList, err := models.DisplayAllItems()
	if err != nil {
		fmt.Println("DisplayAllItems: failed while fetching items detail list with error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed while fetching items list",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"item_list": itemsList,
	})

}

func CurrencyExchange(c *gin.Context) {
	exchangeCurrency := "ZMW"
	exchangeAmount, timeInUnix, err := models.GetExchangeAmountOpenExchange(exchangeCurrency)
	fmt.Println("********amount in conntrollers", exchangeAmount)
	if err != nil {
		timeInUnix = 0
		exchangeAmount = 0.0
		log.Println("CurrencyExchange: failed to get data from openexchangerates.org : ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "failed to get data from api",
			"error":   err,
		})
		return
	}
	countryId, err := models.GetSupportedCountryID(exchangeCurrency)
	if err != nil {
		log.Println("CurrencyExchange: failed while fetching country id: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed while fetching country id",
			"error":   err.Error(),
		})
		return
	}
	//if an error comes while fetching exchange amount, we are not stopping the process, but will use old exchange rate amount
	//as today's exchange rate amount.
	err = models.StoreExchangedAmount(exchangeAmount, timeInUnix, countryId)
	if err != nil {
		log.Println("CurrencyExchange failed while storing data with error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed while storing data with error",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Currency exchange successful",
	})
}

func UpdateAdminSettings(c *gin.Context) {
	profitPerc := c.PostForm("profit")
	profit, err := strconv.ParseFloat(profitPerc, 64)
	if err != nil {
		fmt.Println("UpdateAdminSettings: Error while converting string to float:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Error while converting string to float",
			"error":   err,
		})
		return
	}
	err = models.UpdateAdminSettings(profit)
	if err != nil {
		fmt.Println("UpdateAdminSettings: Error while updating profit percetage:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "failed while updating profit percentage",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully updated profit percentage"})
}
