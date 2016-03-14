import {createSelector} from 'reselect'

const bookmarkSelector = (state) => state.bookmarks
const searchTagSelector = (state) => state.searchTags
const modalInputSelector = (state) => state.modalInput


const appSelector = createSelector(
  bookmarkSelector,
  searchTagSelector,
  modalInputSelector,
  (bookmarks, searchTags, modalInput) => {
    return {
      bookmarks,
      searchTags,
      modalInput
    }
  }
)

export default appSelector
