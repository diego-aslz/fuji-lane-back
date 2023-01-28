package flentities

// Seed the database
func Seed() error {
	return WithRepository(func(r *Repository) error {
		findOrCreate := [][]interface{}{
			{
				Country{ID: 1},
				&Country{ID: 1, Name: "China"},
			},
			{
				Country{ID: 2},
				&Country{ID: 2, Name: "Hong Kong"},
			},
			{
				Country{ID: 3},
				&Country{ID: 3, Name: "Japan"},
			},
			{
				Country{ID: 4},
				&Country{ID: 4, Name: "Singapore"},
			},
			{
				Country{ID: 5},
				&Country{ID: 5, Name: "Vietnam"},
			},
			{
				City{ID: 101},
				&City{ID: 101, CountryID: 1, Name: "Beijing", Latitude: 39.9390731, Longitude: 116.11728},
			},
			{
				City{ID: 102},
				&City{ID: 102, CountryID: 1, Name: "Chengdu", Latitude: 30.6587488, Longitude: 103.935463},
			},
			{
				City{ID: 103},
				&City{ID: 103, CountryID: 1, Name: "Chongqing", Latitude: 29.5551377, Longitude: 106.4084703},
			},
			{
				City{ID: 104},
				&City{ID: 104, CountryID: 1, Name: "Dongguan", Latitude: 22.9764535, Longitude: 113.654243},
			},
			{
				City{ID: 105},
				&City{ID: 105, CountryID: 1, Name: "Guangzhou", Latitude: 23.1259819, Longitude: 112.9476602},
			},
			{
				City{ID: 106},
				&City{ID: 106, CountryID: 1, Name: "Shanghai", Latitude: 31.2246325, Longitude: 121.1965709},
			},
			{
				City{ID: 107},
				&City{ID: 107, CountryID: 1, Name: "Shenyang", Latitude: 41.8055019, Longitude: 123.2964156},
			},
			{
				City{ID: 108},
				&City{ID: 108, CountryID: 1, Name: "Shenzhen", Latitude: 22.5554167, Longitude: 113.913795},
			},
			{
				City{ID: 109},
				&City{ID: 109, CountryID: 1, Name: "Tianjin", Latitude: 39.1252291, Longitude: 117.015353},
			},
			{
				City{ID: 110},
				&City{ID: 110, CountryID: 1, Name: "Wuhan", Latitude: 30.5683366, Longitude: 114.1603012},
			},
			{
				City{ID: 201},
				&City{ID: 201, CountryID: 2, Name: "Hong Kong", Latitude: 22.284736, Longitude: 114.1414606},
			},
			{
				City{ID: 301},
				&City{ID: 301, CountryID: 3, Name: "Fukuoka", Latitude: 33.625038, Longitude: 130.0258401},
			},
			{
				City{ID: 302},
				&City{ID: 302, CountryID: 3, Name: "Kawasaki", Latitude: 35.5562073, Longitude: 139.5723855},
			},
			{
				City{ID: 303},
				&City{ID: 303, CountryID: 3, Name: "Kobe", Latitude: 34.6943656, Longitude: 135.1556806},
			},
			{
				City{ID: 304},
				&City{ID: 304, CountryID: 3, Name: "Kyoto", Latitude: 35.0061653, Longitude: 135.7259306},
			},
			{
				City{ID: 305},
				&City{ID: 305, CountryID: 3, Name: "Nagoya", Latitude: 35.1680838, Longitude: 136.8940904},
			},
			{
				City{ID: 306},
				&City{ID: 306, CountryID: 3, Name: "Osaka", Latitude: 34.69374, Longitude: 135.50218},
			},
			{
				City{ID: 307},
				&City{ID: 307, CountryID: 3, Name: "Saitama", Latitude: 35.915717, Longitude: 139.5787164},
			},
			{
				City{ID: 308},
				&City{ID: 308, CountryID: 3, Name: "Sapporo", Latitude: 43.0595074, Longitude: 141.3354807},
			},
			{
				City{ID: 309},
				&City{ID: 309, CountryID: 3, Name: "Tokyo", Latitude: 35.6735408, Longitude: 139.570305},
			},
			{
				City{ID: 310},
				&City{ID: 310, CountryID: 3, Name: "Yokohama", Latitude: 35.4620149, Longitude: 139.5842306},
			},
			{
				City{ID: 401},
				&City{ID: 401, CountryID: 4, Name: "Singapore", Latitude: 1.3439166, Longitude: 103.7540049},
			},
			{
				City{ID: 501},
				&City{ID: 501, CountryID: 5, Name: "Ho Chi Minh", Latitude: 10.7659164, Longitude: 106.4034602},
			},
			{
				City{ID: 502},
				&City{ID: 502, CountryID: 5, Name: "Hanoi", Latitude: 20.9740874, Longitude: 105.3724915},
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
