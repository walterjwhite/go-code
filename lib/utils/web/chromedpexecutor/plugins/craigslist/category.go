package craigslist

import (
	"fmt"
	"github.com/chromedp/chromedp"
)

type Category int

const (
	Antiques Category = iota + 1
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
	return []chromedp.Action{
	}
}
