// Функция для получения города из localStorage
function getCachedCity() {
    return localStorage.getItem("selectedCity") || "msk";
}


// Функция для сохранения города в localStorage
function setCachedCity(city) {
    localStorage.setItem("selectedCity", city);
}


// Вспомогательная функция для названия города
function getCityName(cityKey) {
    const cities = {
        msk: "Москве",
        spb: "Санкт-Петербурге",
        ekb: "Екатеринбурге"
    };
    return cities[cityKey];
}


document.querySelectorAll(".back-button").forEach(button => {
    button.addEventListener("click", () => {
        // Скрываем все страницы
        document.getElementById("loading-page").classList.add("hidden");
        document.getElementById("results-page").classList.add("hidden");
        document.getElementById("not-found-page")?.classList.add("hidden");

        loadPopularEvents(getCachedCity())
        // Показываем стартовую страницу
        document.getElementById("start-page").classList.remove("hidden");
    });
});


function loadPopularEvents(city) {
    fetch(`/last/${city}`)
        .then(response => response.json())
        .then(data => {
            const title = document.getElementById("popular-title");
            const list = document.getElementById("popular-list");

            // Обновляем заголовок
            const cityName = getCityName(city);
            title.textContent = `Что ищут люди в ${cityName}:`;

            // Очищаем старые данные
            list.innerHTML = "";

            console.log("Недавние мероприятия: ", data);

            // Добавляем новые элементы
            data.forEach((event, index) => {
                if (event.name !== "") {
                    const li = document.createElement("li");

                    const a = document.createElement("a");
                    a.textContent = event.name;
                    a.style.cursor = "pointer";
                    a.style.color = "#3498db";
                    a.style.textDecoration = "none";

                    // При клике запускаем поиск
                    a.addEventListener("click", () => {
                        const searchInput = document.getElementById("search-input");
                        searchInput.value = event.name;

                        // Вызываем обработчик поиска
                        handleSearch();
                    });

                    li.appendChild(a);
                    list.appendChild(li);
                }
            });
        })
        .catch(error => {
            console.error("Ошибка при загрузке популярных событий:", error);
            document.getElementById("popular-title").textContent = "";
            document.getElementById("popular-list").innerHTML = "";
        });
}


function handleSearch() {
    const eventName = document.getElementById("search-input").value.trim();
    const city = document.getElementById("city-select").value;
    const errorDiv = document.getElementById("error-message");

    if (!eventName) {
        errorDiv.classList.remove("hidden");
        return;
    } else {
        errorDiv.classList.add("hidden");
    }

    // Показываем лоадер
    document.getElementById("start-page").classList.add("hidden");
    document.getElementById("loading-page").classList.remove("hidden");

    fetch(`/results/${city}/${encodeURIComponent(eventName)}`)
        .then(response => response.json())
        .then(data => {
            const resultsList = document.getElementById("results-list");
            resultsList.innerHTML = "";

            if (data.length === 0) {
                document.getElementById("loading-page").classList.add("hidden");
                document.getElementById("results-page").classList.add("hidden");
                document.getElementById("not-found-page").classList.remove("hidden");
            } else {
                // Собираем все цены со всех билетов всех мероприятий
                const allPrices = data.flatMap(event => 
                    event.tickets.map(ticket => parseFloat(ticket.price))
                );

                // Находим минимальную цену среди всех билетов
                const globalMinPrice = Math.min(...allPrices);

                data.forEach(event => {
                    // Создаём контейнер для события
                    const eventContainer = document.createElement("div");
                    eventContainer.className = "result-item";

                    // Название мероприятия
                    const eventNameEl = document.createElement("div");
                    eventNameEl.className = "ticket-name";
                    eventNameEl.textContent = event.name;

                    eventContainer.appendChild(eventNameEl);

                    // Сортируем билеты по location, чтобы можно было группировать
                    const ticketsByLocation = {};
                    event.tickets.forEach(ticket => {
                        const loc = ticket.location || "Не указано";
                        if (!ticketsByLocation[loc]) {
                            ticketsByLocation[loc] = [];
                        }
                        ticketsByLocation[loc].push(ticket);
                    });

                    // Выводим билеты по группам
                    for (const location in ticketsByLocation) {
                        const locationHeader = document.createElement("div");
                        locationHeader.className = "ticket-location";
                        locationHeader.textContent = location;

                        eventContainer.appendChild(locationHeader);

                        ticketsByLocation[location].forEach(ticket => {
                            const ticketLink = document.createElement("a");
                            ticketLink.href = ticket.url;
                            ticketLink.target = "_blank";
                            ticketLink.className = "ticket-details";

                            // Если цена билета равна глобальной минимальной, применяем стиль
                            const priceClass = parseFloat(ticket.price) === globalMinPrice ? "lowest-price" : "";
                            ticketLink.innerHTML = `
                                <span>${ticket.date}</span>
                                <span>${ticket.time || "—"}</span>
                                <span class="${priceClass}">${ticket.price}₽</span>
                            `;

                            eventContainer.appendChild(ticketLink);
                        });
                    }

                    resultsList.appendChild(eventContainer);
                });

                document.getElementById("loading-page").classList.add("hidden");
                document.getElementById("results-page").classList.remove("hidden");
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert("Не удалось получить данные.");
            document.getElementById("loading-page").classList.add("hidden");
            document.getElementById("start-page").classList.remove("hidden");
        });
}

document.getElementById("search-button").addEventListener("click", handleSearch);


document.addEventListener("DOMContentLoaded", function () {
    const citySelect = document.getElementById("city-select");
    const cachedCity = getCachedCity();

    // Устанавливаем выбранный город
    citySelect.value = cachedCity;

    // Загружаем популярные события
    loadPopularEvents(cachedCity);

    // При смене города — обновляем кэш и список популярных
    citySelect.addEventListener("change", function () {
        const selectedCity = this.value;
        setCachedCity(selectedCity);
        loadPopularEvents(selectedCity);
    });
});