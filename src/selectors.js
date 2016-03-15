import {createSelector} from 'reselect'

const bookmarkSelector = (state) => state.bookmarks
const searchTagSelector = (state) => state.searchTags.tags
const modalInputSelector = (state) => state.modalInput
const userSelector = (state) => state.user
const inputSelector = (state) => state.searchTags.value


const appSelector = createSelector(
  bookmarkSelector,
  searchTagSelector,
  modalInputSelector,
  userSelector,
  inputSelector,
  (bookmarks, searchTags, modalInput, user, searchInput) => {
    return {
      bookmarks,
      searchTags,
      modalInput,
      user,
      searchInput
    }
  }
)

export default appSelector
