package core

var keyspaceStat [4]map[string]int

func UpdateDbSate(num int, metric string, value int) {
	keyspaceStat[num][metric] = value
}
