package domain

type Country struct {
	ID                int32  `db:"id"`
	CountryName       string `db:"country_name"`
	CountryCode       string `db:"country_code"`
	ContinentName     string `db:"continent_name"`
	CountryNameArabic string `db:"country_name_arabic"`
	CountryPhoneCode  string `db:"country_phone_code"`
	Timezone          string `db:"timezone"`
	UtcOffset         string `db:"utc_offset"`
	Currency          string `db:"currency"`
}
