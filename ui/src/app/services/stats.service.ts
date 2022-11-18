import { Injectable } from '@angular/core';
import { environment } from '@lox/environments/environment';
import { HttpClient } from '@angular/common/http';
import { map, Observable, tap } from 'rxjs';
import { StatsRecordModel } from '@lox/models/stats-record.model';
import { IStatsFlowResponse } from '../interfaces/stats-flow-response.interface';
import { IStatsTotalResponse } from '../interfaces/stats-total-response.interface';
import { IStatsOverviewResponse } from '../interfaces/stats-overview-response.interface';

@Injectable({
  providedIn: 'root'
})
export class StatsService {
  private _records: StatsRecordModel = new StatsRecordModel();
  private _rawData: any;

  constructor(private readonly _http: HttpClient) {}

  public statsTotal(): Observable<any> {
    const url = this.cleanUpUrl(`${this.serviceUrl}/${environment.statsTotalService.endPoint}`);

    return this._http
      .request<IStatsTotalResponse>(environment.statsTotalService.method ?? 'GET', url, {
        headers: {
          'Content-Type': environment.statsTotalService.header.contentType
        }
      })
      .pipe(
        tap((data: IStatsTotalResponse) => {
          this._rawData = data;
        }),
        tap((data: IStatsTotalResponse) => {
          this._records.addTotalRecord(data);
        }),
        map(() => this._records)
      );
  }

  public statsOverview(): Observable<any> {
    const url = this.cleanUpUrl(`${this.serviceUrl}/${environment.statsOverviewService.endPoint}`);

    return this._http
      .request<IStatsOverviewResponse>(environment.statsOverviewService.method ?? 'GET', url, {
        headers: {
          'Content-Type': environment.statsOverviewService.header.contentType
        }
      })
      .pipe(
        tap((data: IStatsOverviewResponse) => {
          this._records.addOverview(data);
        }),
        map(() => this._records)
      );
  }

  public statsFlow(): Observable<any> {
    const url = this.cleanUpUrl(`${this.serviceUrl}/${environment.statsFlowService.endPoint}`);

    return this._http
      .request<IStatsFlowResponse>(environment.statsFlowService.method ?? 'GET', url, {
        headers: {
          'Content-Type': environment.statsFlowService.header.contentType
        }
      })
      .pipe(
        tap((data: IStatsFlowResponse) => {
          this._records.addFlowRecord(data);
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

  public get rawData(): any {
    return this._rawData;
  }
}
