import { IStatsRecord } from '@lox/interfaces/stats-record.interface';
import { StatsService } from '@lox/services/stats.service';
import moment, { Moment } from 'moment';
import { BehaviorSubject, ReplaySubject } from 'rxjs';
import { IPods } from '../interfaces/pod.interface';
import { IStatsFlowResponse } from '../interfaces/stats-flow-response.interface';
import { IStatsTotalResponse } from '../interfaces/stats-total-response.interface';

export class StatsRecordModel {
  private _recordCount: number = 0;
  private _podCount: number = 0;

  private _podNames: string[] = [];
  private _pods: IPods = {};
  private _lastUpdate: ReplaySubject<Moment> = new ReplaySubject<Moment>();
  private _updateFilter: ReplaySubject<void> = new ReplaySubject<void>();
  private _record: IStatsRecord = {};
  private _filteredRecord: IStatsRecord = {};
  private _timeStampArray: Date[] = [];

  constructor() {}

  public addTotalRecord(rawData: IStatsTotalResponse) {
    //key value pair to array
    const data = Object.entries(rawData);
    this._podNames = Object.keys(rawData);
    this._podCount = this._podNames.length;

    // merge new data from rawData into _pods
    Object.keys(rawData).forEach((key: string) => {
      if (!!this._pods[key]) {
        this._pods[key] = { ...this._pods[key], ...rawData[key] };
      } else {
        const recordDummy = Array(this._recordCount).fill({
          packetsSum: 0,
          transmitBytes: 0,
          receivedBytes: 0,
          unknownBytes: 0,
          timeStamp: new Date()
        });
        this._pods[key] = { ...rawData[key], records: recordDummy };
      }
    });
  }

  public addFlowRecord(rawData: IStatsFlowResponse) {
    Object.keys(this._pods).forEach((key: string) => {
      if (!!rawData[key]) {
        this._pods[key].records.push({ ...rawData[key], timeStamp: new Date() });
      } else {
        this._pods[key].records.push({
          packetsSum: 0,
          transmitBytes: 0,
          receivedBytes: 0,
          unknownBytes: 0,
          localTransmitBytes: 0,
          localReceivedBytes: 0,
          timeStamp: new Date()
        });
      }
    });
    this._recordCount++;

    console.log(this._pods);
  }

  public selectPod(pod: string): void {
    this._filteredRecord[pod] = this._record[pod];

    this._updateFilter.next();
  }

  public deSelectPod(pod: string): void {
    delete this._filteredRecord[pod];
    this._updateFilter.next();
  }

  public clearSelection(): void {
    Object.keys(this._filteredRecord).forEach((pod: string) => {
      delete this._filteredRecord[pod];
    });
    this._updateFilter.next();
  }

  get podNames(): string[] {
    return this._podNames;
  }

  get selectedpodNames(): string[] {
    return Object.keys(this._filteredRecord);
  }

  get pods(): IStatsRecord {
    return Object.keys(this._filteredRecord).length > 0 ? this._filteredRecord : this._record;
  }

  get podList(): IPods {
    return this._pods;
  }

  get lastUpdate(): ReplaySubject<Moment> {
    return this._lastUpdate;
  }

  get updateFilter(): ReplaySubject<void> {
    return this._updateFilter;
  }
}
