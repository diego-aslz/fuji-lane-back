package flentities

import "github.com/jinzhu/gorm"

// Seed the database
func Seed() error {
	return WithRepository(func(r *Repository) error {
		findOrCreate := [][]interface{}{
			[]interface{}{
				Country{Model: gorm.Model{ID: 1}},
				&Country{Model: gorm.Model{ID: 1}, Name: "China"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 2}},
				&Country{Model: gorm.Model{ID: 2}, Name: "Hong Kong"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 3}},
				&Country{Model: gorm.Model{ID: 3}, Name: "Japan"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 4}},
				&Country{Model: gorm.Model{ID: 4}, Name: "Singapore"},
			},
			[]interface{}{
				Country{Model: gorm.Model{ID: 5}},
				&Country{Model: gorm.Model{ID: 5}, Name: "Vietnam"},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 101}},
				&City{Model: gorm.Model{ID: 101}, CountryID: 1, Name: "Beijing", Latitude: 39.9390731, Longitude: 116.11728},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 102}},
				&City{Model: gorm.Model{ID: 102}, CountryID: 1, Name: "Chengdu", Latitude: 30.6587488, Longitude: 103.935463},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 103}},
				&City{Model: gorm.Model{ID: 103}, CountryID: 1, Name: "Chongqing", Latitude: 29.5551377, Longitude: 106.4084703},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 104}},
				&City{Model: gorm.Model{ID: 104}, CountryID: 1, Name: "Dongguan", Latitude: 22.9764535, Longitude: 113.654243},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 105}},
				&City{Model: gorm.Model{ID: 105}, CountryID: 1, Name: "Guangzhou", Latitude: 23.1259819, Longitude: 112.9476602},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 106}},
				&City{Model: gorm.Model{ID: 106}, CountryID: 1, Name: "Shanghai", Latitude: 31.2246325, Longitude: 121.1965709},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 107}},
				&City{Model: gorm.Model{ID: 107}, CountryID: 1, Name: "Shenyang", Latitude: 41.8055019, Longitude: 123.2964156},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 108}},
				&City{Model: gorm.Model{ID: 108}, CountryID: 1, Name: "Shenzhen", Latitude: 22.5554167, Longitude: 113.913795},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 109}},
				&City{Model: gorm.Model{ID: 109}, CountryID: 1, Name: "Tianjin", Latitude: 39.1252291, Longitude: 117.015353},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 110}},
				&City{Model: gorm.Model{ID: 110}, CountryID: 1, Name: "Wuhan", Latitude: 30.5683366, Longitude: 114.1603012},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 201}},
				&City{Model: gorm.Model{ID: 201}, CountryID: 2, Name: "Hong Kong", Latitude: 22.284736, Longitude: 114.1414606},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 301}},
				&City{Model: gorm.Model{ID: 301}, CountryID: 3, Name: "Fukuoka", Latitude: 33.625038, Longitude: 130.0258401},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 302}},
				&City{Model: gorm.Model{ID: 302}, CountryID: 3, Name: "Kawasaki", Latitude: 35.5562073, Longitude: 139.5723855},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 303}},
				&City{Model: gorm.Model{ID: 303}, CountryID: 3, Name: "Kobe", Latitude: 34.6943656, Longitude: 135.1556806},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 304}},
				&City{Model: gorm.Model{ID: 304}, CountryID: 3, Name: "Kyoto", Latitude: 35.0061653, Longitude: 135.7259306},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 305}},
				&City{Model: gorm.Model{ID: 305}, CountryID: 3, Name: "Nagoya", Latitude: 35.1680838, Longitude: 136.8940904},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 306}},
				&City{Model: gorm.Model{ID: 306}, CountryID: 3, Name: "Osaka", Latitude: 34.69374, Longitude: 135.50218},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 307}},
				&City{Model: gorm.Model{ID: 307}, CountryID: 3, Name: "Saitama", Latitude: 35.915717, Longitude: 139.5787164},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 308}},
				&City{Model: gorm.Model{ID: 308}, CountryID: 3, Name: "Sapporo", Latitude: 43.0595074, Longitude: 141.3354807},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 309}},
				&City{Model: gorm.Model{ID: 309}, CountryID: 3, Name: "Tokyo", Latitude: 35.6735408, Longitude: 139.570305},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 310}},
				&City{Model: gorm.Model{ID: 310}, CountryID: 3, Name: "Yokohama", Latitude: 35.4620149, Longitude: 139.5842306},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 401}},
				&City{Model: gorm.Model{ID: 401}, CountryID: 4, Name: "Singapore", Latitude: 1.3439166, Longitude: 103.7540049},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 501}},
				&City{Model: gorm.Model{ID: 501}, CountryID: 5, Name: "Ho Chi Minh", Latitude: 10.7659164, Longitude: 106.4034602},
			},
			[]interface{}{
				City{Model: gorm.Model{ID: 502}},
				&City{Model: gorm.Model{ID: 502}, CountryID: 5, Name: "Hanoi", Latitude: 20.9740874, Longitude: 105.3724915},
			},
		}

		for _, pairs := range findOrCreate {
			if err := r.Where(pairs[0]).FirstOrCreate(pairs[1]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
