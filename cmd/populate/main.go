package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/icrowley/fake"
	"github.com/nerde/fuji-lane-back/flconfig"
)

const userPasswords = "123456789"

var propertyImageURLs, unitImageURLs []string
var cities []*flentities.City

func main() {
	flconfig.LoadConfiguration()

	if flconfig.Config.Stage == "production" {
		log.Fatalf("Command \"populate\" should not be run in %s stage!", flconfig.Config.Stage)
	}

	if err := flentities.Reset(); err != nil {
		log.Fatal(err.Error())
	}

	if err := flentities.Seed(); err != nil {
		log.Fatal(err.Error())
	}

	propertyImageURLs = []string{
		"https://2qibqm39xjt6q46gf1rwo2g1-wpengine.netdna-ssl.com/wp-content/uploads/2018/06/12184802_web1_M1-Alderwood-EDH-180611.jpg",
		"https://upload.wikimedia.org/wikipedia/commons/thumb/1/14/Forsyth_Barr_Building%2C_Christchurch_02.JPG/240px-Forsyth_Barr_Building%2C_Christchurch_02.JPG",
		"https://triblive.com/csp/mediapool/sites/dt.common.streams.StreamServer.cls?STREAMOID=k7NPPpUTmn_wjH7LMbdfgM$daE2N3K4ZzOUsqbU5sYsqXrmU32razg_4hODQsDzEWCsjLu883Ygn4B49Lvm9bPe2QeMKQdVeZmXF$9l$4uCZ8QDXhaHEp3rvzXRJFdy0KqPHLoMevcTLo3h8xh70Y6N_U_CryOsw6FTOdKL_jpQ-&amp;CONTENTTYPE=image/jpeg",
		"https://www.gannett-cdn.com/presto/2018/08/24/PWES/13a9c900-41ad-4333-a5fb-f4e93e369785-GATEWAY_PHOTO.jpg?width=534&height=712&fit=bounds&auto=webp",
	}

	unitImageURLs = []string{
		"https://37b3a77d7df28c23c767-53afc51582ca456b5a70c93dcc61a759.ssl.cf3.rackcdn.com/1024x768/54850_3971_001.jpg",
		"https://cdngeneral.rentcafe.com/dmslivecafe/3/622878/slider_The-Gallery-Apartments-1020.jpg?quality=85&scale=both&",
		"https://g5-assets-cld-res.cloudinary.com/image/upload/q_auto,f_auto,c_fill,g_center,h_1100,w_2000/v1504090383/g5/g5-c-i7yxybw5-mission-rock-single/g5-cl-iap27qrg-mountain-view-apartment-homes/uploads/interior-4.jpg",
		"https://www.onni.com/wp-content/uploads/2016/11/Rental-Apartment-Page-new-min.jpg",
		"https://cdn.vox-cdn.com/thumbor/E0jNRUTI81RBBRMSA_ZU7vq7I4g=/0x0:2400x1602/1200x675/filters:focal(682x772:1066x1156)/cdn.vox-cdn.com/uploads/chorus_image/image/54241701/LINEA_NicholasJamesPhoto_8.0.jpeg",
	}

	if err := listCities(); err != nil {
		log.Fatal(err.Error())
	}

	createAccount()
}

func createAccount() error {
	city := randomCity()
	return flentities.WithRepository(func(r *flentities.Repository) error {
		acc := &flentities.Account{
			Name:      fake.Company(),
			Phone:     str(fake.Phone()),
			CountryID: ui(city.CountryID),
		}

		if err := r.Save(acc).Error; err != nil {
			return err
		}

		owner := &flentities.User{
			AccountID: ui(acc.ID),
			Name:      str(fake.FullName()),
			Email:     fake.EmailAddress(),
		}

		owner.SetPassword(userPasswords)

		if err := r.Save(owner).Error; err != nil {
			return err
		}

		for i := 0; i < 10; i++ {
			publishedAt := time.Now()
			publishedAt = publishedAt.AddDate(0, 0, randomInt(365)*-1)
			prop := &flentities.Property{
				Name:            str(fake.ProductName()),
				PublishedAt:     &publishedAt,
				EverPublished:   true,
				AccountID:       acc.ID,
				Address1:        str(fake.StreetAddress()),
				Address2:        str(fake.StreetAddress()),
				Address3:        str(fake.StreetAddress()),
				PostalCode:      str(fake.Zip()),
				CityID:          ui(city.ID),
				CountryID:       ui(city.CountryID),
				Latitude:        city.Latitude + ((float32(randomInt(400)) - 200.0) / 1000.0),
				Longitude:       city.Longitude + ((float32(randomInt(400)) - 200.0) / 1000.0),
				MinimumStay:     in(randomInt(4) + 1),
				NearestAirport:  str(fake.MaleFullName()),
				NearestSubway:   str(fake.Street()),
				NearbyLocations: str(fake.MaleFullName()),
				Overview:        str(randomOverview()),
			}

			if err := r.Save(prop).Error; err != nil {
				return err
			}

			for _, img := range buildImages(propertyImageURLs, 3) {
				img.Property = prop

				if err := r.Save(img).Error; err != nil {
					return err
				}
			}

			for _, am := range buildAmenities(flentities.PropertyAmenityTypes) {
				am.PropertyID = ui(prop.ID)

				c := 0
				if r.Where(am).Table("amenities").Count(&c); c > 0 {
					continue
				}

				if err := r.Save(am).Error; err != nil {
					return err
				}
			}

			for j := 0; j < 6; j++ {
				publishedAt := time.Now()
				publishedAt = publishedAt.AddDate(0, 0, randomInt(365)*-1)

				unit := &flentities.Unit{
					PublishedAt:   &publishedAt,
					EverPublished: true,
					Property:      prop,
					Name:          fake.FullName(),
					Overview:      str(randomOverview()),
					Bedrooms:      randomInt(8),
					Bathrooms:     randomInt(10),
					SizeM2:        randomInt(100) + 20,
					MaxOccupancy:  in(randomInt(12) + 2),
					Count:         randomInt(50) + 5,
				}

				if err := r.Save(unit).Error; err != nil {
					return err
				}

				for _, img := range buildImages(unitImageURLs, 3) {
					img.Unit = unit

					if err := r.Save(img).Error; err != nil {
						return err
					}
				}

				for _, am := range buildAmenities(flentities.UnitAmenityTypes) {
					am.UnitID = ui(unit.ID)

					c := 0
					if r.Where(am).Table("amenities").Count(&c); c > 0 {
						continue
					}

					if err := r.Save(am).Error; err != nil {
						return err
					}
				}

				minNights := []int{1, 2, 7, 30, 90, 180, 360}
				prices := []int{}
				for idx, nights := range minNights {
					variation := 30
					prices = append(prices, nights*(randomInt(variation)+130-(idx*variation/2)))
				}
				pricesCount := randomInt(7) + 1
				for j := 0; j < pricesCount; j++ {
					pr := &flentities.Price{
						Unit:      unit,
						MinNights: minNights[j],
						Cents:     prices[j] * 100,
					}

					if err := r.Save(pr).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func str(s string) *string {
	return &s
}

func ui(u uint) *uint {
	return &u
}

func in(i int) *int {
	return &i
}

func randomCity() *flentities.City {
	return cities[randomInt(len(cities))]

}

var random = rand.New(rand.NewSource(time.Now().Unix()))

func randomInt(limit int) int {
	return random.Intn(limit)
}

func buildAmenities(types []*flentities.AmenityType) []*flentities.Amenity {
	amenities := []*flentities.Amenity{}

	l := len(types)
	for i := 0; i < l/2; i++ {
		typ := flentities.AmenityTypes[randomInt(l)]
		amenities = append(amenities, &flentities.Amenity{Type: typ.Code, Name: str(typ.Name)})
	}

	return amenities
}

func buildImages(urls []string, count int) []*flentities.Image {
	images := []*flentities.Image{}

	for j := 0; j < count; j++ {
		images = append(images, &flentities.Image{
			URL:      urls[randomInt(len(urls))],
			Name:     fake.ProductName(),
			Type:     "img/jpeg",
			Size:     randomInt(1024 * 1024 * (1 + randomInt(5))),
			Uploaded: true,
			Position: j + 1,
		})
	}

	return images
}

func randomOverview() string {
	return fmt.Sprintf("<p>%s</p>", fake.Paragraph())
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func listCities() error {
	cities = []*flentities.City{}
	return flentities.WithRepository(func(r *flentities.Repository) error {
		return r.Find(&cities).Error
	})
}
