// Lightweight database statistics tracked alongside key operations.
package core

var keyspaceStat [4]map[string]int

// UpdateDbSate sets a metric value for the selected logical database.
func UpdateDbSate(num int, metric string, value int) {
	keyspaceStat[num][metric] = value
}
