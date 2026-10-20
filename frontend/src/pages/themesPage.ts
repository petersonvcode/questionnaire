import { hideElement, Page, showElement } from "./page";
import { countQuestionsWithTags, CountQuestionsWithTags, getTags, Tag } from "../questions";

export class ThemesQuestionPage implements Page {
  private readonly pageElement: HTMLElement;
  private readonly themesContainerElement: HTMLElement;
  private readonly themesListElement: HTMLElement;
  private readonly qAvailableElement: HTMLElement;
  private readonly qCountElements: HTMLButtonElement[];
  private readonly confirmTagsButtonElement: HTMLButtonElement;

  /** All available tags */
  private allTags: Tag[] = [];
  /** Tags selected by the user */
  private selectedTags: Tag[] = [];
  /** Given the selected tags, there are this many questions available, also, what other tags can be used to further filter the questions */
  private questionsCount: CountQuestionsWithTags | null = null;
  /** How many questions the user wants to answer */
  private desiredQuestionsCount: number | null = null;

  constructor(
    /** Callback to be called when user confirms the tags and the desired questions count */
    onConfirmTags: (tags: Tag[], count: number) => void
  ) {
    const pageElement = document.getElementById('themes-question-page');
    if (!pageElement)
      throw new Error('Page element not found');
    this.pageElement = pageElement;

    const themesContainerElement = document.getElementById('themes-container');
    if (!themesContainerElement)
      throw new Error('Themes container element not found');
    this.themesContainerElement = themesContainerElement;

    const themesListElement = document.getElementById('themes-list');
    if (!themesListElement)
      throw new Error('Themes list element not found');
    this.themesListElement = themesListElement;

    const qAvailableElement = document.getElementById('q-available');
    if (!qAvailableElement)
      throw new Error('Q available element not found');
    this.qAvailableElement = qAvailableElement;

    const confirmTagsButtonElement = document.getElementById('confirm-tags-btn');
    if (!confirmTagsButtonElement || !(confirmTagsButtonElement instanceof HTMLButtonElement))
      throw new Error('Confirm tags button element not found');
    this.confirmTagsButtonElement = confirmTagsButtonElement;
    this.confirmTagsButtonElement.addEventListener('click', () => onConfirmTags(this.selectedTags, this.desiredQuestionsCount ?? 0));

    const qCountBtns = document.getElementsByClassName('q-count-btn');
    this.qCountElements = Array.from(qCountBtns) as HTMLButtonElement[];
    for (const btn of this.qCountElements)
      btn.addEventListener('click', () => {
        this.changeDesiredQuestionsCount(parseInt(btn.textContent || '0'));
      });
  }

  show = () => {
    this.themesListElement.innerHTML = ''
    showElement(this.pageElement)
    showElement(this.themesContainerElement)
    countQuestionsWithTags(this.selectedTags)
      .then(questionsCount => {
        this.questionsCount = questionsCount;
        this.qAvailableElement.textContent = `Total disponíveis: ${questionsCount.totalCount}`;
      })
      .catch(error => console.error('Error counting questions with tags', error));
    getTags().then(tags => {
      this.allTags = tags;
      this.renderTagList();
    })
  }

  hide = () => {
    hideElement(this.pageElement)
    hideElement(this.themesContainerElement)
    this.selectedTags = [];
    this.desiredQuestionsCount = null;
    this.questionsCount = null;
    this.qAvailableElement.textContent = 'Disponíveis: 0';
    for (const btn of this.qCountElements) {
      btn.disabled = true;
      btn.classList.remove('selected');
    }
    this.disableConfirmButton()
  }

  private renderTagList = () => {
    this.themesListElement.innerHTML = '';
    const selectedIDs = new Set(this.selectedTags.map(t => t.id));
    const anySelected = selectedIDs.size > 0;
    const otherAvailableTags = this.questionsCount?.otherTagsCount ?? [];
    const otherAvailableTagIDs = new Set(otherAvailableTags.map(t => t.tagID));

    const tagsToRender = anySelected
      ? this.allTags.filter(t => otherAvailableTagIDs.has(t.id) || selectedIDs.has(t.id))
      : this.allTags;

    for (const tag of tagsToRender) {
      const isSelected = selectedIDs.has(tag.id);
      const tagElement = document.createElement('div');
      tagElement.classList.add('theme-item');
      if (isSelected)
        tagElement.classList.add('selected');
      tagElement.textContent = tag.tag;
      tagElement.addEventListener('click', () => this.onTagClick(tag));
      this.themesListElement.appendChild(tagElement);
    }
  }

  private onTagClick = async (tag: Tag) => {
    const isSelected = this.selectedTags.some(s => s.id === tag.id);
    if (isSelected)
      this.selectedTags = this.selectedTags.filter(s => s.id !== tag.id);
    else
      this.selectedTags.push(tag);

    this.qAvailableElement.textContent = '...';

    this.questionsCount = await countQuestionsWithTags(this.selectedTags);
    this.qAvailableElement.textContent = `Total disponíveis: ${this.questionsCount?.totalCount ?? 0}`;

    this.renderTagList();

    this.updateDesiredCountButtons();

    if (
      this.desiredQuestionsCount !== null &&
      this.questionsCount !== null &&
      this.desiredQuestionsCount > 0 &&
      this.questionsCount.totalCount >= this.desiredQuestionsCount
    )
      this.enableConfirmButton()
    else
      this.disableConfirmButton()
  }

  private enableConfirmButton = () => {
    this.confirmTagsButtonElement.disabled = false;
    this.confirmTagsButtonElement.style.display = 'block';
  }

  private disableConfirmButton = () => {
    this.confirmTagsButtonElement.disabled = true;
    this.confirmTagsButtonElement.style.display = 'block';
  }

  private changeDesiredQuestionsCount = (newCount: number) => {
    if (newCount > (this.questionsCount?.totalCount || 0)) {
      console.warn('Cannot set desired questions count to a value greater than the available questions count');
      return
    }
    this.desiredQuestionsCount = newCount;

    this.updateDesiredCountButtons()

    if (this.desiredQuestionsCount !== null && this.questionsCount?.totalCount !== null && this.desiredQuestionsCount > 0)
      this.enableConfirmButton()
    else
      this.disableConfirmButton()
  }

  private updateDesiredCountButtons = () => {
    for (const btn of this.qCountElements) {
      const btnValue = parseInt(btn.textContent || '0');
      if (btnValue > (this.questionsCount?.totalCount || 0)) {
        btn.disabled = true;
        if (btn.classList.contains('selected')) {
          btn.classList.remove('selected');
          this.desiredQuestionsCount = null;
        }
      } else
        btn.disabled = false;

      if (btnValue === (this.desiredQuestionsCount || 0))
        btn.classList.add('selected');
      else
        btn.classList.remove('selected');
    }
  }
}