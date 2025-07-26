package api

import (
	"errors"
	"fmt"
	"github.com/Roflan4eg/test_work/internal/models"
	"github.com/Roflan4eg/test_work/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	db *storage.SubscriptionStorage
}

func New(storage *storage.SubscriptionStorage) *Handler {
	return &Handler{db: storage}
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	sub := models.Subscription{}
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.Error(errors.New("invalid data"))
		return
	}
	if err := sub.ParseDates(); err != nil {
		c.Error(err)
		return
	}
	if err := h.db.Create(&sub); err != nil {
		c.Error(err)
		return
	}
	sub.FormatDates()
	c.JSON(http.StatusCreated, sub)
}

func (h *Handler) GetSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.New(`filed "id" should be integer`))
		return
	}
	sub, err := h.db.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	sub.FormatDates()
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) ListSubscriptions(c *gin.Context) {
	subs, err := h.db.List()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errors.New(`filed "id" should be integer`))
		return
	}
	sub, err := h.db.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	updateSub := models.CUSubscription{}
	if err = c.ShouldBindJSON(&updateSub); err != nil {
		c.Error(errors.New("invalid data"))
		return
	}
	sub.ServiceName = updateSub.ServiceName
	sub.Price = updateSub.Price
	sub.RawStartDate = updateSub.StartDate
	sub.RawEndDate = updateSub.EndDate
	if err = sub.ParseDates(); err != nil {
		c.Error(err)
		return
	}
	if err = h.db.Update(sub); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(fmt.Errorf("invalid field ID"))
		return
	}
	if err = h.db.Delete(id); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetPriceForPeriod(c *gin.Context) {
	req := models.PeriodRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(fmt.Errorf("invalid data"))
		return
	}
	start, err := time.Parse("01-2006", req.Start)
	if err != nil {
		c.Error(fmt.Errorf("invalid start_date format, expected MM-YYYY"))
	}
	end, err := time.Parse("01-2006", req.End)
	if err != nil {
		c.Error(fmt.Errorf("invalid end_date format, expected MM-YYYY"))
		return
	}
	res, err := h.db.GetSubsForPeriod(start, end, req.UserID, req.ServiceName)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": res})
}
