module.exports = [
  {
    context: ["/api/planets"],
    target: "http://localhost:3000",
    logLevel: "debug",
  },
  {
    context: ["/api/planets/*"],
    target: "http://localhost:3000",
    logLevel: "debug",
  },
  //{
  //  context: ["/api/"],
  //  target: "https://swapi.dev/",
  //  secure: false,
  //  logLevel: "debug",
  //},
];
