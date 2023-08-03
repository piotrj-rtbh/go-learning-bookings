package dbrepo

// we'll create any function that will be available to the interface repository.DatabaseRepo

func (m *postgresDBRepo) AllUsers() bool {
	return true
}
