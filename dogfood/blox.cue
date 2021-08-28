{
  data_dir: "data"
  schemata_dir: "schemata"
  template_dir: "data/tpl"
  static_dir: "static"
  prebuild: [ {
    name: "images"
    executable: "images_impl"
  }]
}
