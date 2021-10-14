//go:generate go run -tags generate ../../generate/listpages/main.go -ListOps=DescribeFleets,DescribeImageBuilders,DescribeStacks -Export=yes
//go:generate go run -tags generate ../../generate/tags/main.go -ListTags=yes -ServiceTagsMap=yes -UpdateTags=yes
// ONLY generate directives and package declaration! Do not add anything else to this file.

package appstream
