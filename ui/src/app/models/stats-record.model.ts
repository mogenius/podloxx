import { IStatsRecord } from '@lox/interfaces/stats-record.interface';
import moment, { Moment } from 'moment';
import { ReplaySubject } from 'rxjs';
import { IPods } from '../interfaces/pod.interface';
import { IStatsFlowResponse } from '../interfaces/stats-flow-response.interface';
import { IStatsOverviewResponse } from '../interfaces/stats-overview-response.interface';
import { IStatsTotalResponse } from '../interfaces/stats-total-response.interface';

export class StatsRecordModel {
  /// RAW DATA
  private _overviewStats: IStatsOverviewResponse;

  /// BASE INFORMATION
  private _recordCount: number = 0;
  private _lastUpdate: ReplaySubject<Moment> = new ReplaySubject<Moment>();

  private _podCount: number = 0;
  private _podNames: string[] = [];

  private _pods: IPods = {};

  private _record: IStatsRecord = {};
  private _filteredRecord: IStatsRecord = {};
  private _sortedNameList: string[] | undefined = undefined;

  private _currentSortKey: string = 'podName';
  private _currentSortDirection: 'ASC' | 'DSC' = 'ASC';

  /// UPDATE EVENTS
  private _updateTotalEvent: ReplaySubject<void> = new ReplaySubject<void>();
  private _updateFlowEvent: ReplaySubject<void> = new ReplaySubject<void>();
  private _updateOverviewEvent: ReplaySubject<void> = new ReplaySubject<void>();

  constructor() {}

  public addTotalRecord(rawData: IStatsTotalResponse) {
    this._podNames = Object.keys(rawData);

    let changed: boolean = false;

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
        this._pods[key] = { ...rawData[key], records: recordDummy, index: this._podCount++ };
        changed = true;
      }
    });
    if (changed) {
      this.sortPods(this._currentSortKey, this._currentSortDirection);
    }
    this._updateTotalEvent.next();
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
    this._updateFlowEvent.next();
  }

  public addOverview(data: IStatsOverviewResponse): void {
    this._overviewStats = data;
    this._updateOverviewEvent.next();
  }

  public selectPod(pod: string): void {
    this._filteredRecord[pod] = this._record[pod];
  }

  public sortPods(key: string, direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._currentSortDirection = direction;
    this._currentSortKey = key;

    switch (key) {
      case 'podName':
        this.sortPodsByName(direction);
        break;
      case 'namespace':
        this.sortPodsByNamespace(direction);
        break;
      case 'externalTraffic':
        this.sortPodsByExternalTraffic(direction);
        break;
      case 'localTraffic':
        this.sortPodsByLocalTraffic(direction);
        break;
      case 'age':
        this.sortPodsByAge(direction);
        break;
      case 'connections':
        this.sortPodsByConnections(direction);
        break;
      default:
        this.sortPodsByName(direction);
        break;
    }
  }

  //ff sort Pods by PodName
  private sortPodsByName(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      if (direction === 'ASC') {
        return a.localeCompare(b);
      } else {
        return b.localeCompare(a);
      }
    });
  }

  //ff sort Pods by Namespace
  private sortPodsByNamespace(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      if (direction === 'ASC') {
        return this._pods[a].namespace.localeCompare(this._pods[b].namespace);
      } else {
        return this._pods[b].namespace.localeCompare(this._pods[a].namespace);
      }
    });
  }

  //ff Sort Pods by external Traffic
  private sortPodsByExternalTraffic(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      // Calculate the sum of transmitted and received traffic of pod a and b
      const aTotal =
        this._pods[a].transmitBytes -
        this._pods[a].localTransmitBytes +
        this._pods[a].receivedBytes -
        this._pods[a].localReceivedBytes;

      const bTotal =
        this._pods[b].transmitBytes -
        this._pods[b].localTransmitBytes +
        this._pods[b].receivedBytes -
        this._pods[b].localReceivedBytes;

      if (direction === 'ASC') {
        return bTotal - aTotal;
      } else {
        return aTotal - bTotal;
      }
    });
  }

  //ff Sort Pods by Local Traffic
  private sortPodsByLocalTraffic(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      // Calculate the sum of transmitted and received traffic of pod a and b
      const aTotal = this._pods[a].localTransmitBytes + this._pods[a].localReceivedBytes;
      const bTotal = this._pods[b].localTransmitBytes + this._pods[b].localReceivedBytes;

      if (direction === 'ASC') {
        return bTotal - aTotal;
      } else {
        return aTotal - bTotal;
      }
    });
  }

  //ff sort pods by age
  private sortPodsByAge(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      if (direction === 'ASC') {
        return moment(this._pods[a].startTime).diff(moment(this._pods[b].startTime));
      } else {
        return moment(this._pods[b].startTime).diff(moment(this._pods[a].startTime));
      }
    });
  }

  //ff sort pods by connections
  private sortPodsByConnections(direction: 'ASC' | 'DSC' = 'ASC'): void {
    this._sortedNameList = Object.keys(this._pods).sort((a: string, b: string) => {
      if (direction === 'ASC') {
        return Object.keys(this._pods[b].connections).length - Object.keys(this._pods[a].connections).length;
      } else {
        return Object.keys(this._pods[a].connections).length - Object.keys(this._pods[b].connections).length;
      }
    });
  }

  //ff Deselect Pod if one was selectet over navigation
  public deSelectPod(pod: string): void {
    delete this._filteredRecord[pod];
  }

  //ff Clear all selected Pods to normal array
  public clearSelection(): void {
    Object.keys(this._filteredRecord).forEach((pod: string) => {
      delete this._filteredRecord[pod];
    });
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

  get overviewStats(): IStatsOverviewResponse {
    return this._overviewStats;
  }

  get sortedNameList(): string[] | undefined {
    if (!!this.selectedpodNames && this.selectedpodNames.length > 0) {
      return this._sortedNameList?.filter((podName) => {
        return this.selectedpodNames.includes(podName);
      });
    } else {
      return this._sortedNameList;
    }
  }

  get updateTotalEvent(): ReplaySubject<void> {
    return this._updateTotalEvent;
  }

  get updateOverviewEvent(): ReplaySubject<void> {
    return this._updateOverviewEvent;
  }

  get updateFlowEvent(): ReplaySubject<void> {
    return this._updateFlowEvent;
  }
}
