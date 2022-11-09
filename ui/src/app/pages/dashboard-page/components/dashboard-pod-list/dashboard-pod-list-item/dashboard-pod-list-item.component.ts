import { Component, Input, OnInit } from '@angular/core';
import { StatsService } from '../../../../../services/stats.service';

@Component({
  selector: '[lox-dashboard-pod-list-item]',
  templateUrl: './dashboard-pod-list-item.component.html',
  styleUrls: ['./dashboard-pod-list-item.component.scss']
})
export class DashboardPodListItemComponent implements OnInit {
  @Input() key: string;

  constructor(private readonly statsService: StatsService) {}

  ngOnInit(): void {
    console.log(this.statsService.records.podList[this.key]);
  }

  get pod(): any {
    return this.statsService.records.podList[this.key];
  }

  get startTime(): number {
    return Math.floor((new Date().getTime() - new Date(this.pod.startTime).getTime()) / (1000 * 60 * 60 * 24));
  }

  // get number of connections
  get connections(): number {
    return Object.keys(this.pod.connections).length;
  }
}
