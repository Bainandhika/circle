package tool

import (
	"circle-fiber/app/config"
	"circle-fiber/lib/model"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/golang/freetype"
)

func printBill(bill model.Bill) string {
	var billString string

	billString = `
=================================================
   ____   ___   ____      ____   ___      ______
  /    | |   | |  _ \    /    | |   |    |      |
 /   __| |   | | | | |  /   __| |   |    |   ___|
|   /    |   | | |_| / |   /    |   |    |  |___
|   |    |  ===[ KEEP IT ROUNDED ]===    |   ___|
|   \__  |   | |     \ |   \__  |   |__  |  |___
 \     | |   | |  |\  \ \     | |      | |      |
  \____| |___| |__| \__\ \____| |______| |______|
=================================================
`

	billString = billString + fmt.Sprintf("Created By: %s\n", bill.CreatedBy)
	billString = billString + fmt.Sprintf("Order Title: %s\n", bill.ResultOrderMain.OrderTitle)
	billString = billString + fmt.Sprintf("Total: %2.f\n", bill.ResultOrderMain.Total)
	billString = billString + fmt.Sprintln("Additional:")
	for _, add := range bill.ResultOrderMain.AdditionalAndDiscounts.Additional {
		billString = billString + fmt.Sprintf("  - %s: %.2f\n", add.Type, add.Cost)
	}
	billString = billString + fmt.Sprintln("Discounts:")
	for _, disc := range bill.ResultOrderMain.AdditionalAndDiscounts.Discounts {
		billString = billString + fmt.Sprintf("  - %s: %.2f\n", disc.Type, disc.Cost)
	}
	billString = billString + fmt.Sprintf("Total Payment: %.2f\n", bill.ResultOrderMain.TotalPayment)
	billString = billString + fmt.Sprintln()

	billString = billString + fmt.Sprintln("Order Details User:")
	billString = billString + fmt.Sprintln()

	for _, person := range bill.ResultOrderUsers {
		billString = billString + fmt.Sprintf("Name: %s\n", person.Name)
		billString = billString + fmt.Sprintln("Items:")
		billString = billString + fmt.Sprintln(strings.Repeat("-", 60))
		billString = billString + fmt.Sprintf("%-20s %-10s %-15s %-10s\n", "Item", "Quantity", "Price Per Item", "Total")
		billString = billString + fmt.Sprintln(strings.Repeat("-", 60))
		for _, item := range person.OrderItems {
			billString = billString + fmt.Sprintf("%-20s %-10d %-15.2f %-10.2f\n", item.Item, item.Quantity, item.PricePerItem, item.Total)
		}
		billString = billString + fmt.Sprintln(strings.Repeat("-", 60))
		billString = billString + fmt.Sprintf("Total: %.2f\n", person.Total)
		billString = billString + fmt.Sprintf("Part of Order: %.2f\n", person.PartOfOrder)
		billString = billString + fmt.Sprintf("Price to Pay: %.2f\n", person.PriceToPay)
		billString = billString + fmt.Sprintln()
	}

	return billString
}

func CreateBillImageFile(bill model.Bill) error {
	billString := printBill(bill)

	// Load TrueType font
	fontBytes, err := os.ReadFile(config.App.FontAssetPath)
	if err != nil {
		return errors.New("Error reading font file: " + err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return errors.New("Error parsing font: " + err.Error())
	}

	// Create a new image with white background
	img := image.NewRGBA(image.Rect(0, 0, 800, 1200))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Set font and font size
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(14)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.Black))

	// Draw bill string on the image
	pt := freetype.Pt(10, 20)
	lines := strings.Split(billString, "\n")
	for _, line := range lines {
		_, err = c.DrawString(line, pt)
		if err != nil {
			return errors.New("Error drawing string: " + err.Error())
		}
		pt.Y += c.PointToFixed(16) // Adjust line height
	}

	// Save the image to a file
	fileName := fmt.Sprintf("%s%s.png", config.App.BillPath, bill.ResultOrderMain.OrderTitle)
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("Error creating file: " + err.Error())
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return errors.New("Error encoding image: " + err.Error())
	}

	return nil
}

func CreateBillTextFile(bill model.Bill) error {
	billString := printBill(bill)

	fileName := fmt.Sprintf("%s%s.txt", config.App.BillPath, bill.ResultOrderMain.OrderTitle)
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("Error creating file: " + err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(billString)
	if err != nil {
		return errors.New("Error writing file: " + err.Error())
	}

	return nil
}
