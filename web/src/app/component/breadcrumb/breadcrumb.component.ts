import {Component, Input, OnInit} from '@angular/core';
import {ActivatedRoute, Router, RouterLink} from "@angular/router";

interface Breadcrumb {
  label: string;
  url: any[];
}

@Component({
  selector: 'layout-partial-app-breadcrumb',
  standalone: true,
  imports: [
    RouterLink
  ],
  templateUrl: './breadcrumb.component.html',
  styleUrl: './breadcrumb.component.scss'
})
export class LayoutPartialBreadcrumbComponent implements OnInit{
  @Input() breadcrumbs: Breadcrumb[] = []

  constructor() {}

  ngOnInit(): void {

  }

  protected genUrl(data:any[]): string{
    return data.join('/')
  }
}
