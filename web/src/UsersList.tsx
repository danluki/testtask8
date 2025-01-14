import React, { useState, useEffect } from 'react';

const UsersList = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(4); 

  useEffect(() => {
    const fetchUsers = async () => {
      setLoading(true);
      const response = await fetch(`https://jsonplaceholder.typicode.com/users?_page=${page}&_limit=3`);
      const data = await response.json();
      setUsers(data);
      setLoading(false);
    };

    fetchUsers();
  }, [page, filter]);

  const handleFilterChange = (e) => {
    setFilter(e.target.value);
    setPage(1); 
  };

  const handleNextPage = () => {
    if (page < totalPages) {
      setPage(page + 1);
    }
  };

  const handlePreviousPage = () => {
    if (page > 1) {
      setPage(page - 1);
    }
  };

  const filteredUsers = users.filter(user => 
    user.name.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div>
      <input
        type="text"
        placeholder="Фильтровать по имени"
        value={filter}
        onChange={handleFilterChange}
      />
      {loading ? (
        <p>Загрузка...</p>
      ) : (
        <ul>
          {filteredUsers.map(user => (
            <li key={user.id}>{user.name}</li>
          ))}
        </ul>
      )}
      <div>
        <button onClick={handlePreviousPage} disabled={page === 1}>Предыдущая</button>
        <button onClick={handleNextPage} disabled={page === totalPages}>Следующая</button>
      </div>
    </div>
  );
};

export default UsersList;