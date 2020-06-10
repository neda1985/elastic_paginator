package pagination

var sample = `{
  "hits":{
    "hits":[
      {"_source":{"t": "1"}},
      {"_source":{"t": "2"}},
      {"_source":{"t": "3"}},
      {"_source":{"t": "4"}},
      {"_source":{"t": "5"}}
    ],
    "total":{
      "value":5
    }
  }
}`

var sample2 = `{
  "hits":{
    "hits":[
 
    ],
    "total":{
      "value":0
    }
  }
}`
