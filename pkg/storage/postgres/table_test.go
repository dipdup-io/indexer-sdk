package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"gopkg.in/yaml.v3"
)

type person struct {
	bun.BaseModel `bun:"person"`

	Id   int64  `bun:"id,pk,autoincrement"`
	Name string `bun:"name,notnull"`
	Age  uint8  `bun:"age,notnull"`
}

// TableName -
func (person) TableName() string {
	return "person"
}

type personTable struct {
	*Table[*person]
}

func newPersonTable(conn *database.Bun) personTable {
	return personTable{
		Table: NewTable[*person](conn),
	}
}

// Flat -
func (p person) Flat() []any {
	return []any{p.Id, p.Name, p.Age}
}

// Columns -
func (person) Columns() []string {
	return []string{"id", "name", "age"}
}

// TableTestSuite -
type TableTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	db            *database.Bun
	table         personTable
}

// SetupSuite -
func (s *TableTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "postgres:15",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	storage, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, func(ctx context.Context, conn *database.Bun) error {
		_, err := conn.DB().NewCreateTable().IfNotExists().Model(new(person)).Exec(ctx)
		return err
	})
	s.Require().NoError(err)
	s.db = storage.db

	s.Require().NoError(err)

	s.table = newPersonTable(s.db)
}

// TearDownSuite -
func (s *TableTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TableTestSuite) TestSave() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	pers := person{
		Name: "Billy",
		Age:  12,
	}

	s.Require().NoError(s.table.Save(ctx, &pers))

	p, err := s.table.GetByID(ctx, uint64(pers.Id))
	s.Require().NoError(err)
	s.Require().Equal(p.Name, pers.Name)
	s.Require().Equal(p.Age, pers.Age)
}

func (s *TableTestSuite) TestUpdate() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	pers := person{
		Name: "Billy",
		Age:  12,
		Id:   1,
	}

	s.Require().NoError(s.table.Update(ctx, &pers))

	p, err := s.table.GetByID(ctx, uint64(pers.Id))
	s.Require().NoError(err)
	s.Require().Equal(p.Name, pers.Name)
	s.Require().Equal(p.Age, pers.Age)
}

func (s *TableTestSuite) TestList() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	list, err := s.table.List(ctx, 2, 1, storage.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(list, 2)
	s.Require().EqualValues(list[0].Id, 2)
	s.Require().EqualValues(list[1].Id, 3)
}

func (s *TableTestSuite) TestCursorList() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	list, err := s.table.CursorList(ctx, 2, 2, storage.SortOrderAsc, storage.ComparatorGt)
	s.Require().NoError(err)
	s.Require().Len(list, 2)
	s.Require().EqualValues(list[0].Id, 3)
	s.Require().EqualValues(list[1].Id, 4)
}

func (s *TableTestSuite) TestGetById() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	pers, err := s.table.GetByID(ctx, 2)
	s.Require().NoError(err)
	s.Require().EqualValues(pers.Id, 2)
	s.Require().EqualValues(pers.Name, "Anna")
	s.Require().EqualValues(pers.Age, 55)
}

func (s *TableTestSuite) TestIsNoRows() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	_, err = s.table.GetByID(ctx, 11)
	s.Require().Error(err)
	s.Require().Equal(s.table.IsNoRows(err), true)
}

func (s *TableTestSuite) TestLastId() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	lastId, err := s.table.LastID(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(lastId, 4)
}

func (s *TableTestSuite) TestTransaction() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	transactable := NewTransactable(s.db)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := transactable.BeginTransaction(ctx)
	s.Require().NoError(err)

	newPerson := person{
		Name: "Billy",
		Age:  12,
	}
	updPerson := person{
		Name: "Villy",
		Age:  22,
		Id:   1,
	}

	s.Require().NoError(tx.Add(ctx, &newPerson))
	s.Require().NoError(tx.Update(ctx, &updPerson))
	s.Require().NoError(tx.Flush(ctx))
	s.Require().NoError(tx.Close(ctx))

	newP, err := s.table.GetByID(ctx, uint64(newPerson.Id))
	s.Require().NoError(err)
	s.Require().EqualValues(newP.Id, newPerson.Id)
	s.Require().Equal(newP.Name, newPerson.Name)
	s.Require().EqualValues(newP.Age, newPerson.Age)

	updP, err := s.table.GetByID(ctx, uint64(updPerson.Id))
	s.Require().NoError(err)
	s.Require().EqualValues(updP.Id, updPerson.Id)
	s.Require().Equal(updP.Name, updPerson.Name)
	s.Require().EqualValues(updP.Age, updPerson.Age)
}

func (s *TableTestSuite) TestTransactionCopyFrom() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixtures"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	transactable := NewTransactable(s.db)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer ctxCancel()

	tx, err := transactable.BeginTransaction(ctx)
	s.Require().NoError(err)

	res, err := tx.Exec(ctx, "delete from person")
	s.Require().NoError(err)
	s.Require().EqualValues(4, res)

	f, err := os.Open("test/fixtures/person.yml")
	s.Require().NoError(err)
	defer f.Close()

	var ps []person
	err = yaml.NewDecoder(f).Decode(&ps)
	s.Require().NoError(err)

	data := make([]storage.Copiable, len(ps))
	for i := range ps {
		data[i] = ps[i]
	}

	err = tx.CopyFrom(ctx, ps[0].TableName(), data)
	s.Require().NoError(err)

	s.Require().NoError(tx.Rollback(ctx))
	s.Require().NoError(tx.Close(ctx))

	count, err := s.db.DB().NewSelect().Model(&person{}).Count(ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(4, count)

}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}
