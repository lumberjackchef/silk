package main

type ComponentMeta struct {
  ProjectName   string  `json:"project_name"`
  ComponentName string  `json:"component_name"`
  InitDate      string  `json:"init_date"`
  Version       string  `json:"version"`
  Description   string  `json:"description"`
}
