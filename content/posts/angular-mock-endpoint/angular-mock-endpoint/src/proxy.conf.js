const PROXY_CONFIG = [
  {
    context: ["/api/"],
    target: "https://swapi.dev/",
    secure: false,
  },
];

module.exports = PROXY_CONFIG;
