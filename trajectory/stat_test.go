package trajectory

import "testing"

func TestParseStatBundled(t *testing.T) {
	request := "{\"request\": \"requestid.path.appid.version.module.instanceid\"," +
		"\"items\": {\"cpu_usage\": \"0.0|g\", \"memory_usage\": \"0.0|g\"}}"

	stats := ParseStat([]byte(request))

	compareStat1 := Stat{
		"requestid",
		"requestid.path.appid.version.module.instanceid.cpu_usage",
		"0.0",
		"g",
	}
	verifyStat(t, request, stats[0], compareStat1)

	compareStat2 := Stat{
		"requestid",
		"requestid.path.appid.version.module.instanceid.memory_usage",
		"0.0",
		"g",
	}
	verifyStat(t, request, stats[1], compareStat2)
}

func verifyStat(t *testing.T, request string, stat, compareStat Stat) {
	if stat.Parent != compareStat.Parent {
		t.Errorf("ParseStat(%v).Parent = %v, want %v", request, stat.Parent,
			compareStat.Parent)
	}
	if stat.Id != compareStat.Id {
		t.Errorf("ParseStat(%v).Id = %v, want %v", request, stat.Id,
			compareStat.Id)
	}
	if stat.Value != compareStat.Value {
		t.Errorf("ParseStat(%v).Value = %v, want %v", request, stat.Value,
			compareStat.Value)
	}
	if stat.Type != compareStat.Type {
		t.Errorf("ParseStat(%v).Type = %v, want %v", request, stat.Type,
			compareStat.Type)
	}
}

func TestParseSingleStatNoBundle(t *testing.T) {
	request := "requestid.path.appid.version.module.instanceid.cpu_usage:0.0|g"

	stats := ParseStat([]byte(request))

	compareStat1 := Stat{
		"requestid",
		"requestid.path.appid.version.module.instanceid.cpu_usage",
		"0.0",
		"g",
	}
	verifyStat(t, request, stats[0], compareStat1)
}

func TestParseMultipleStatNoBundle(t *testing.T) {
	request := "requestid.path.appid.version.module.instanceid.cpu_usage:0.0|g" +
		",requestid.path.appid.version.module.instanceid.memory_usage:2.5|s"

	stats := ParseStat([]byte(request))

	compareStat1 := Stat{
		"requestid",
		"requestid.path.appid.version.module.instanceid.cpu_usage",
		"0.0",
		"g",
	}
	verifyStat(t, request, stats[0], compareStat1)

	compareStat2 := Stat{
		"requestid",
		"requestid.path.appid.version.module.instanceid.memory_usage",
		"2.5",
		"s",
	}
	verifyStat(t, request, stats[1], compareStat2)
}
