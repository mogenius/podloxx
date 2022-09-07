import { Moment } from 'moment';

export interface IStatsRecord {
  [key: string]: {
    index: number;
    timeStamps: Date[];
    receiveBytes: number[];
    transmitBytes: number[];
  };
}
