import { Injectable } from '@angular/core';
import { environment } from '@lox/environments/environment';
import { HttpClient } from '@angular/common/http';
import { map, Observable, Subject, tap } from 'rxjs';
import { StatsRecordModel } from '@lox/models/stats-record.model';
import { IStatsResponse } from '@lox/interfaces/stats-response.interface';

@Injectable({
  providedIn: 'root'
})
export class StatsService {
  private _records: StatsRecordModel = new StatsRecordModel();
  private _rawData: string;

  constructor(private readonly _http: HttpClient) {}

  public stats(): Observable<any> {
    const url = this.cleanUpUrl(`${this.serviceUrl}/${environment.statsTotalService.endPoint}`);

    return this._http
      .request<any>(environment.statsTotalService.method ?? 'GET', url, {
        headers: {
          'Content-Type': environment.statsTotalService.header.contentType
        }
      })
      .pipe(
        tap((data: IStatsResponse[]) => {
          this._records.addRecord(data);
        }),
        map(() => this._records)
      );
  }

  private cleanUpUrl(str: string): string {
    return str.replace(/([^:]\/)\/+/g, '$1');
  }

  private get serviceUrl(): string {
    return this.cleanUpUrl(`${environment.baseUrl}/`);
  }

  public get records(): StatsRecordModel {
    return this._records;
  }
}
