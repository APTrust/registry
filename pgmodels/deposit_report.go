package pgmodels

import ()

// TODO: Filters...

// DepositReport contains info about member deposits and the costs
// of those deposits. This struct does not implement the usual pgmodel
// interface, nor does it map to a single underlying table or view.
// This struct merely represents to the output of a reporting query.
type DepositReport struct {
	InstitutionName string
	FileCount       int64
	ObjectCount     int64
	TotalSize       int64
	StorageOption   string
}

// TODO: Pull in costs from StorageOptions and add method
//       to calculate cost line by line.
//       Or roll this into the query.

/*
Base working query for this report:

select
  i."name" as instituition_name,
  count(gf.id) as file_count,
  count(distinct(gf.intellectual_object_id)) as object_count,
  sum(gf.size) as total_size,
  gf.storage_option
from generic_files gf
left join institutions i on i.id = gf.institution_id
where gf.state = 'A' and gf.created_at < '2021-10-01'
group by cube (i."name", gf.storage_option)
order by i.name, gf.storage_option

*/
