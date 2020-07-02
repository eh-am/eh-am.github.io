import { Component, OnInit } from "@angular/core";
import {
  PlanetsDataService,
  Planet,
  DetailedPlanet,
} from "./planets-data.service";
import { map } from "rxjs/operators";
import { Observable } from "rxjs";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
})
export class AppComponent implements OnInit {
  constructor(private planetsDataService: PlanetsDataService) {}
  planets$: Observable<Planet[] | null>;
  detailedPlanet$: Observable<DetailedPlanet> | null;
  detailedPlanet: DetailedPlanet | null = null;

  ngOnInit() {
    this.planets$ = this.planetsDataService
      .getPaginated()
      .pipe(map((response) => response.results));
  }

  isSelected(p: Planet) {
    if (!this.detailedPlanet) {
      return false;
    }

    // May not be the best way to validate uniqueness
    // But it's fine for our case
    return this.detailedPlanet.name === p.name;
  }

  onClick(p: Planet) {
    this.planetsDataService.getDetails(p).subscribe((dp) => {
      this.detailedPlanet = dp;
    });
  }
}
