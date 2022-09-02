import { Moment } from 'moment';

export interface IStatsRecord {
  [key: string]: {
    index: number;
    receiveBytes: number[];
    transmitBytes: number[];
  };
}
