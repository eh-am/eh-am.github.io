import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";

@Injectable({
  providedIn: "root",
})
export class PlanetsDataService {
  constructor(private http: HttpClient) {}

  getAll() {
    return this.http.get("/api/planets");
  }
}
