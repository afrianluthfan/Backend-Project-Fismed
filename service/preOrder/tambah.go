package preOrder

import (
	"backend_project_fismed/constanta"
	"backend_project_fismed/model"
	"backend_project_fismed/utility"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Posting(c *gin.Context) {
	var input model.PurchaseOrder
	//var res model.PurchaseOrder
	var id int

	if c.GetHeader("content-type") == "application/x-www-form-urlencoded" || c.GetHeader("content-type") == "application/x-www-form-urlencoded; charset=utf-8" {

		if err := c.Bind(&input); err != nil {
			utility.ResponseError(c, "Input Data Tidak Berhasil !")
		}

	} else {

		if err := c.BindJSON(&input); err != nil {
			utility.ResponseError(c, "Input Data Tidak Berhasil !")
		}

	}

	ctx := context.Background()
	tx, err := DBConnect.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(err.Error())
	}

	log.Println("Data Input :", input)

	query := `
		INSERT INTO purchase_order (
			nama_suplier, 
			nomor_po, 
			tanggal, 
			catatan_po, 
			prepared_by, 
			prepared_jabatan, 
			approved_by, 
			approved_jabatan, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by, 
			sub_total, 
			pajak, 
			total
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)RETURNING id`

	// Menjalankan query insert
	err = tx.QueryRow(context.Background(), query, input.NamaSuplier, input.NomorPO, input.Tanggal,
		input.CatatanPO, input.PreparedBy, input.PreparedJabatan, input.ApprovedBy, input.ApprovedJabatan,
		time.Now(), "ADMIN", time.Now(), "ADMIN", input.SubTotal, input.Pajak, input.Total).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		utility.ResponseError(c, constanta.ErrQuery1)
		return
	}

	if len(input.Item) > 0 {
		for i, item := range input.Item {
			QueryItem := `
				INSERT INTO item_buyer (
					po_id, 
					name, 
					quantity, 
					price, 
					discount, 
					amount
				) VALUES ($1, $2, $3, $4, $5, $6)`

			_, err = tx.Exec(context.Background(), QueryItem, id, item.Name, item.Quantity, item.Price, item.Discount, item.Amount)
			if err != nil {
				tx.Rollback(ctx)
				utility.ResponseError(c, constanta.ErrQuery2)
				return
			}

			log.Println("Data Item ke : ", i, " Berhasil di Input!")

		}
	}

	if err := tx.Commit(ctx); err != nil {
		utility.ResponseError(c, constanta.ErrCommit)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data added successfully", "status": true})
}

func Inquiry(c *gin.Context) {
	var input model.PurchaseOrder
	var res model.PurchaseOrder

	if c.GetHeader("content-type") == "application/x-www-form-urlencoded" || c.GetHeader("content-type") == "application/x-www-form-urlencoded; charset=utf-8" {

		if err := c.Bind(&input); err != nil {
			utility.ResponseError(c, "Input Data Tidak Berhasil !")
		}

	} else {

		if err := c.BindJSON(&input); err != nil {
			utility.ResponseError(c, "Input Data Tidak Berhasil !")
		}

	}

	log.Println("Data Input :", input)

	res.Item = make([]model.ItemBuyer, len(input.Item))

	var Subtotal int

	if len(input.Item) > 0 {
		for i, item := range input.Item {
			log.Println("Item Barang Ke :", i)

			res.Item[i].Name = item.Name
			res.Item[i].Quantity = item.Quantity
			res.Item[i].Price = "Rp. " + utility.FormatRupiah(item.Price)
			res.Item[i].Discount = item.Discount

			QuantitiInt, err := strconv.Atoi(item.Quantity)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid Quantity format"})
				return
			}

			HargaSatuanInt, err := strconv.Atoi(item.Price)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid Harga Satuan format"})
				return
			}

			DiscountInt, err := strconv.Atoi(item.Discount)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid Discount format"})
				return
			}

			subtotalperitem := QuantitiInt*HargaSatuanInt - (QuantitiInt * HargaSatuanInt * DiscountInt / 100)

			Subtotal = Subtotal + subtotalperitem

			subtotalperitemstring := strconv.Itoa(subtotalperitem)

			res.Item[i].Amount = "Rp. " + utility.FormatRupiah(subtotalperitemstring)
		}

		pajak := Subtotal * 11 / 100
		total := pajak + Subtotal

		res.Total = "Rp. " + utility.FormatRupiah(strconv.Itoa(total))
		res.Pajak = "Rp. " + utility.FormatRupiah(strconv.Itoa(pajak))
		res.SubTotal = "Rp. " + utility.FormatRupiah(strconv.Itoa(Subtotal))

		currentTime := time.Now()
		timeString := currentTime.Format("2006-01-02")

		res.NamaSuplier = input.NamaSuplier
		res.NomorPO = utility.GenerateNomorPO()
		res.Tanggal = timeString
		res.CatatanPO = input.CatatanPO
		res.PreparedBy = input.PreparedBy
		res.PreparedJabatan = input.PreparedJabatan
		res.ApprovedBy = input.ApprovedBy
		res.ApprovedJabatan = input.ApprovedJabatan

		c.JSON(http.StatusOK, gin.H{"message": "Inquiry Purcase Order Success !", "data": res, "status": true})

	} else {
		utility.ResponseError(c, "Item Empty!")
	}

}
