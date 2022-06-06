# OWID-API 🌎 

This project is a Proof of Concept which implements a GraphQL API for exploring [**OurWorldInData** (OWID) datasets](https://github.com/owid/owid-datasets).  

[OWID](https://ourworldindata.org/) is a scientific online publication that focuses on large global problems such as poverty, disease, hunger, climate change, war, existential risks, and inequality.
[The goal of OWID](https://ourworldindata.org/about) is to make the knowledge on the big problems accessible and understandable. As they say on their homepage, Our World in Data’s mission is to publish the “research and data to make progress against the world’s largest problems”.

## How does OWID-API work?

The data in [OWID](https://github.com/owid/owid-datasets) is stored as CSV files. 
Every dataset has schema information (see example [here](https://github.com/owid/owid-datasets/blob/master/datasets/Time%20use%20in%20Sweden%20-%20Statistics%20Sweden/datapackage.json)). 
Every dataset has 2 keys to access to the different variables. 
These 2 keys are **entity** and **year**. The first might be whatever dimension of the information like country, person names and so on.  

OWID-API implements a [GraphQL Schema](https://github.com/margostino/owid-api/blob/master/graph/schema.graphql) which is generated from the metadata which and it fetches the information from the CSV files. 

The Schema and Server need to be re-generated when there is a new update in [the datasets](https://github.com/owid/owid-datasets).
Currently, this is a manual process. There are 2 step to re-generate:

1. Generate Schema
2. Generate Server  

You can get the current schema executing the following queries:  

1. Get Full Schema. See [here]()
2. Get All Dataset names. See [here]()  



### Generate Schema

There is one Go function to generate the schema. This process read all files `datapackage.json` from [this repo](https://github.com/owid/owid-datasets/tree/master/datasets) and create a new file `schema.graphql`. 

```bash
> make schema.gen 
```

### Generate Server

Finally the following command will generate the server using [gqlgen](https://github.com/99designs/gqlgen).

```bash
> make server.gen
```

## Usage

There are 2 ways to explore OWID datasets: Rest API or GraphiQL Playground.

### Rest API:  

**POST** https://owid-api.vercel.app/api/query  
```json
{
  "query": "{\n\ttime_use_in_sweden_statistics_sweden(entity: \"Gainful employment\", year:1990){ time_allocation_weekday_women \n\t}}"
}
```

For now, no authentication required. 

### Playground:

Endpoint: https://owid-api.vercel.app/api/playground

**For example**: query 2 datasets, __time_use_in_sweden_statistics_sweden__ AND  __time_use_in_finland_statistics_finland__

```graphql
{
  time_use_in_sweden_statistics_sweden(entity: "Gainful employment", year: 1990) {
    time_allocation_weekday_women
  }
  time_use_in_finland_statistics_finland(entity: "Free time", year: 1987) {
    time_allocation_all_statistics_finland
    time_allocation_women_statistics_finland
  }
}
```

## Response

Either Rest API or Playground query, the response is always a GraphQL representation. See [here](https://github.com/graphql/graphql-spec/blob/main/spec/Section%207%20--%20Response.md) for more details.  

## Schema Design

### Naming
When generating schema, naming dataset and variables is using the same normalizer which applies the same rules, like `toLower` or `replace(...)`.
This is an initial approach and in next iterations it should be improved in order to handle shorter naming and more descriptive since for example:
`total_value_of_exports_by_country_to_world_percgdp_owid_calculations_based_on_fouquin_and_hugot_cepii_2016_and_other_sources` is clearly not really easy to remember or deal with.  

### Types
Another decision was the types of each variable. A first data inspection shows that all variables, except __entity__ are numbers, and since some of them are decimals, the type for all variables is `Float`.

### Arguments

This project assumes that every dataset has 2 arguments: `Entity:String` and `Year:Int`. 