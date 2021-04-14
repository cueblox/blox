import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';

function SearchBar({ filterResults }) {
  const [searchTerm, setSearchTerm] = useState('');

  function handleSearch(term) {
    setSearchTerm(term)
    filterResults(term)
  }

  return (
    <div className="w-full px-4 sm:px-6 my-4 sm:my-8 max-w-4xl relative mx-auto text-gray-500 text-center w-full sm:w-1/2">
      <input className="border-2 border-gray-300 bg-white h-10 px-5 pr-16 rounded-lg text-sm focus:outline-none w-full"
        onChange={(event) => handleSearch(event.target.value)} type="search" name="search" placeholder="Search" />
      <button type="submit" className="absolute right-0 top-0 mr-10 mt-2 focus:outline-none">
        <FontAwesomeIcon icon={faSearch} className="w-5 text-gray-400" />
      </button>
    </div>
  )
}

export default SearchBar;
