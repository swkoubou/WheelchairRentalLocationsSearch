package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/swkoubou/Hack_U_App/Models"
	"os"
	"unicode/utf8"
)

// IDatabaseServiceの実装
type Mysql struct {
	db *sql.DB
}

const decimalNumber = 10

func NewMysql() *Mysql {
	mysqlUser := os.Getenv("HACK_U_App_MYSQL_USER")
	if utf8.RuneCountInString(mysqlUser) == 0 {
		panic("環境変数が読み込めませんでした:" + mysqlUser)
	}
	mysqlPassword := os.Getenv("HACK_U_App_MYSQL_PASSWORD")
	if utf8.RuneCountInString(mysqlPassword) == 0 {
		panic("環境変数が読み込めませんでした:" + mysqlPassword)
	}
	mysqlHost := os.Getenv("HACK_U_App_MYSQL_HOST")
	if utf8.RuneCountInString(mysqlPassword) == 0 {
		panic("環境変数が読み込めませんでした:" + mysqlHost)
	}
	src := fmt.Sprintf("%s:%s@tcp(%s:3306)/hack_u_db", mysqlUser, mysqlPassword, mysqlHost)
	db, err := sql.Open("mysql", src)
	if err != nil {
		panic(err)
	}
	return &Mysql{db: db}
}

// 北西 南東
func (mysql *Mysql) FindAllLocation() (locations []*Models.WheelchairRentalLocation, err error) {
	query := `SELECT location_id, name, x(location), y(location), address_supplement, address, phone_number, email, web_site_url
	FROM wheelchair_rental_Locations;`
	stmt, err := mysql.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		location := Models.NewWheelchairRentalLocation()
		err := rows.Scan(
			&location.LocationId,
			&location.Name,
			&location.Location.Lat,
			&location.Location.Lng,
			&location.AddressSupplement,
			&location.Address,
			&location.PhoneNumber,
			&location.Email,
			&location.WebSiteUrl)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (mysql *Mysql) FindOneLocation(locationId uint64) (location *Models.WheelchairRentalLocation, err error) {
	query := `SELECT location_id, name, x(location), y(location), address_supplement, address, phone_number, email, web_site_url
FROM wheelchair_rental_Locations WHERE Location_id = ?;`
	stmt, err := mysql.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(locationId)
	defer rows.Close()
	location = Models.NewWheelchairRentalLocation()
	err = stmt.QueryRow(locationId).Scan(
		&location.LocationId,
		&location.Name,
		&location.Location.Lat,
		&location.Location.Lng,
		&location.AddressSupplement,
		&location.Address,
		&location.PhoneNumber,
		&location.Email,
		&location.WebSiteUrl)
	if err != nil {
		return nil, err
	}
	return location, nil

}


func (mysql *Mysql) FindAllTags() (tags []*Models.Tag, err error) {
	query := `SELECT * FROM tag`
	stmt, err := mysql.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		tag := Models.NewTag()
		err := rows.Scan(
			&tag.TagId,
			&tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (mysql *Mysql) FindTag(location *Models.WheelchairRentalLocation) (tags []*Models.Tag, err error) {
	query := `SELECT tag.tag_id, tag_name FROM wheelchair_rental_Locations_tag
         JOIN tag ON tag.tag_id = wheelchair_rental_Locations_tag.tag_id
WHERE Location_id = ?;`
	stmt, err := mysql.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(location.LocationId)
	defer rows.Close()
	for rows.Next() {
		tag := Models.NewTag()
		err := rows.Scan(
			&tag.TagId,
			&tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}