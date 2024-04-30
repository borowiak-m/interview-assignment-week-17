package handlers

// /GET records from Mongo
// 1. request payload:
// ---JSON: 4 fields
// startDate date
// endDate date
// minCount int
// maxCount int

// 2. response payload:
// ---JSON: 3 fields
// code: status code (0 for success)
// msg: status description (success)
// records: array of filtered items including fields
// - key
// - createdAt
// - totalCount = sum of count in the document/table
