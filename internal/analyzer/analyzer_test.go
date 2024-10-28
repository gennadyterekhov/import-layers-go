package analyzer

//func TestCanDetectConst(t *testing.T) {
//	results := analysistest.Run(t, analysistest.TestData(), New().Analyzer)
//
//	assert.Equal(t, 1, len(results))
//	assert.Nil(t, results[0].Err)
//	assert.Equal(t, 0, len(results[0].Diagnostics))
//}
//
//func TestReportIfNoConst(t *testing.T) {
//	results := analysistest.Run(t, analysistest.TestData(), New().Analyzer, "noconst")
//
//	assert.Equal(t, 1, len(results))
//	assert.Equal(t, 1, len(results[0].Diagnostics))
//}

// TODO fix
//func TestLevels(t *testing.T) {
//	reportsDownlevel(t)
//	//doesNotReportUplevel(t)
//}

//func reportsDownlevel(t *testing.T) {
//	var result analysistest.Result
//	fmt.Println("reportsDownlevel before run")
//
//	results := analysistest.Run(t, analysistest.TestData(), New().Analyzer, "low", "high")
//	fmt.Println("reportsDownlevel after run")
//
//	require.Greater(t, len(results), 0)
//	result = *results[0]
//
//	require.Greater(t, len(result.Diagnostics), 0)
//	diag := result.Diagnostics[0]
//
//	assert.Equal(t, "wrong level cannot import low from high", diag.Message)
//}
//
//func doesNotReportUplevel(t *testing.T) {
//	var result analysistest.Result
//	results := analysistest.Run(t, analysistest.TestData(), New().Analyzer, "./uplevel/low", "./uplevel/high")
//
//	assert.Equal(t, 1, len(results))
//	result = *results[0]
//	assert.NoError(t, result.Err)
//
//	assert.Equal(t, 0, len(result.Diagnostics))
//}
