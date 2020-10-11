package craigslist

import (
	"fmt"
	"github.com/chromedp/chromedp"
)

type Category int

const (
	//*[@id="new-edit"]/div/label/label[1]/div/span
	Antiques Category = iota + 1
	//*[@id="new-edit"]/div/label/label[2]/div/span
	Appliances
	ArtsAndCrafts
	AtvsUtvsSnowmobiles
	AutoParts
	AutoWheelsAndTires
	Avaiation
	BabyAndKid
	Barter
	BicycleParts
	Bicycles
	BoatParts
	Boats
	BooksAndMagazines
	BusinessCommercial
	CarsAndTrucks
	CdsDvdsVhs
	CellPhones
	ClothingAndAccessories
	Collectibles
	ComputerParts
	Computers
	Electronics
	FarmAndGarden
	FreeStuff
	Furniture
	GarageAndMovingSales
	GeneralForSale
	HealthAndBeauty
	HeavyEquipment
	HouseholdItems
	Jewelry
	Materials
	MotorcycleParts
	MotorcyclesScooters
	MusicalInstruments
	PhotoVideo
	Rvs
	SportingGoods
	Tickets
	Tools
	ToysAndGames
	Trailers
	VideoGaming
	Wanted
)

func (p *CraigslistPost) doCategory() []chromedp.Action {
	// for sale by owner
	return []chromedp.Action{
		// bicycle parts - by owner
		chromedp.Click(fmt.Sprintf("//*[@id=\"new-edit\"]/div/label/label[%v]", p.Category)),
	}
}
