import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-homepage-page',
  imports: [CommonModule, RouterLink],
  templateUrl: './homepage-page.component.html',
  styleUrl: './homepage-page.component.scss',
})
export class HomepagePageComponent {}
