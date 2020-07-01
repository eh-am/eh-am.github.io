import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";

export interface StarWarsApiPaginatedResult<T> {
  count: number;
  next: string | null;
  previous: string | null;
  results: T[];
}

export interface Planet {
  name: string;
  url: string;
}

export interface DetailedPlanet extends Planet {
  climate: string;
  diameter: string;
  terrain: string;
  gravity: string;
}

@Injectable({
  providedIn: "root",
})
export class PlanetsDataService {
  constructor(private http: HttpClient) {}

  getPaginated() {
    return this.http.get<StarWarsApiPaginatedResult<Planet>>("/api/planets");
  }

  getDetails(p: Planet) {
    return this.http.get<DetailedPlanet>(p.url);
  }
}
