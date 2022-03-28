package db

import (
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"github/Outsider565/gostudy/apiservice/utils"
	"os"
	"strings"
)

func (con *Connection) addConstraint(tableName, constraintName, checkClause string) error {
	_, err := con.db.Model().Exec(fmt.Sprintf("SELECT create_constraint_if_not_exists('%s','%s','ALTER TABLE %s ADD CONSTRAINT %s CHECK(%s);') ", tableName, constraintName, tableName, constraintName, checkClause))
	if err != nil {
		return err
	}
	return nil
}
func (con *Connection) addForeignKey(tableName, fkColumns, parentTable, parentColumns string) error {
	constraintName := "fk_" + strings.ReplaceAll(fkColumns, ",", "_")
	_, err := con.db.Model().Exec(fmt.Sprintf("SELECT create_constraint_if_not_exists('%s','%s','ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY(%s) REFERENCES %s(%s);') ", tableName, constraintName, tableName, constraintName, fkColumns, parentTable, parentColumns))
	if err != nil {
		return err
	}
	return nil
}
func (con *Connection) CreateSchema() error {
	db := con.db
	for _, model := range allModels {
		if err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		}); err != nil {
			return err
		}
	}
	//ref: https://stackoverflow.com/questions/6801919/postgres-add-constraint-if-it-doesnt-already-exist
	fStr := `create or replace function create_constraint_if_not_exists (
    t_name text, c_name text, constraint_sql text
	) 
	returns void AS
	$$
	begin
		-- Look for our constraint
		if not exists (select conname
					   from pg_constraint 
					   where conname=c_name) then
			execute constraint_sql;
		end if;
	end;
	$$ language 'plpgsql'`
	if _, err := db.Model().Exec(fStr); err != nil {
		return err
	}
	err := con.addConstraint("teach_ins", "chk_index", "class_index>=1 and class_index<=14")
	if err != nil {
		return err
	}
	err = con.addConstraint("users", "chk_gender", "gender>=0 and gender<=2")
	if err != nil {
		return err
	}
	if _, err := db.Model().Exec("CREATE EXTENSION if not exists pg_trgm;"); err != nil {
		return err
	}
	err = con.addConstraint("courses", "chk_semester", "semester IN ($$Spring$$,$$Autumn$$,$$Summer$$)")
	if err != nil {
		return err
	}
	err = con.addForeignKey("locate_ins", "building_no", "buildings", "building_no")
	if err != nil {
		return err
	}
	err = con.addForeignKey("locate_ins", "classroom_no", "classrooms", "classroom_no")
	if err != nil {
		return err
	}
	err = con.addForeignKey("teach_ins", "classroom_no", "classrooms", "classroom_no")
	if err != nil {
		return err
	}
	err = con.addForeignKey("teach_ins", "course_id,year,semester", "courses", "course_id,year,semester")
	if err != nil {
		return err
	}
	err = con.addForeignKey("studies", "user_id", "users", "user_id")
	if err != nil {
		return err
	}
	err = con.addForeignKey("studies", "classroom_no", "classrooms", "classroom_no")
	if err != nil {
		return err
	}
	err = con.addForeignKey("takes", "user_id", "users", "user_id")
	if err != nil {
		return err
	}
	err = con.addForeignKey("takes", "course_id,year,semester", "courses", "course_id,year,semester")
	if err != nil {
		return err
	}
	return nil
}
func (con *Connection) DropAllTables() error {
	db := con.db
	for _, model := range allModels {
		if err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		}); err != nil {
			return err
		}
	}
	return nil
}

var buildings = [...]Building{
	{
		"H2", "二教", 31.297664, 121.504083,
	},
	{
		BuildingNo:   "H2",
		BuildingName: "",
		Latitude:     31.297664,
		Longitude:    121.504083,
	},
	{
		BuildingNo:   "H3NoH30",
		BuildingName: "三教",
		Latitude:     31.298067,
		Longitude:    121.504265,
	},
	{
		BuildingNo:   "H4",
		BuildingName: "四教",
		Latitude:     31.29847,
		Longitude:    121.50146,
	},
	{
		BuildingNo:   "H5",
		BuildingName: "五教",
		Latitude:     31.295564,
		Longitude:    121.504764,
	},
	{
		BuildingNo:   "H6",
		BuildingName: "六教",
		Latitude:     31.294918,
		Longitude:    121.504212,
	},
	{
		BuildingNo:   "HGX",
		BuildingName: "光华楼西辅楼",
		Latitude:     31.305438,
		Longitude:    121.511061,
	},
	{
		BuildingNo:   "HGD",
		BuildingName: "光华楼东辅楼",
		Latitude:     31.305804,
		Longitude:    121.512350,
	},
	{
		BuildingNo:   "HQ",
		BuildingName: "新闻学院",
		Latitude:     31.305268,
		Longitude:    121.516262,
	},
	{
		BuildingNo:   "J",
		BuildingName: "江湾",
		Latitude:     31.34272,
		Longitude:    121.513724,
	},
	{
		BuildingNo:   "Z",
		BuildingName: "张江",
		Latitude:     31.196204,
		Longitude:    121.60448,
	},
	{
		BuildingNo:   "F",
		BuildingName: "枫林",
		Latitude:     31.203596,
		Longitude:    121.458046,
	},
}

func (con *Connection) InitBuildingInfo() error {
	db := con.db
	for _, building := range buildings {
		_, err := db.Model(&building).OnConflict("DO NOTHING").Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

var TestUser = User{
	UserId:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	OpenId:    "None",
	AvatarUrl: "test_url",
	NickName:  "test_nickname",
	Gender:    0,
	City:      "test_city",
	Province:  "test_province",
	Country:   "test_country",
}

func (con *Connection) InitAdmin() error {
	Adminer := Admin{
		AdminId:      os.Getenv("ADMIN_ID"),
		PasswordHash: utils.Hashstr(os.Getenv("ADMIN_PASSWD")),
	}
	_, err := con.db.Model(&Adminer).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	_, err = con.db.Model(&TestUser).OnConflict("DO NOTHING").Insert()
	return err
}
