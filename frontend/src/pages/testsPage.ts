import { Page, showElement, hideElement } from "./page";

export class TestPage implements Page {
  private readonly pageElement: HTMLElement;
  constructor() {
    const pageElement = document.getElementById('test-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;
  }

  show = () => showElement(this.pageElement)

  hide = () => hideElement(this.pageElement)
}