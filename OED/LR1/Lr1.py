import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import scipy.stats as stats

# 1. Загрузка данных (GDP per capita по странам за 2022 год)
url = "https://databank.worldbank.org/data/download/GDP.csv"
df = pd.read_csv(url, skiprows=4)

# Выбираем данные за 2022 год и удаляем пропущенные значения
df = df[['Country Name', '2022']].dropna()
df.columns = ['Country', 'GDP_per_capita']
data = df['GDP_per_capita'].values

# 2. Основные статистические характеристики
mean = np.mean(data)
median = np.median(data)
mode = stats.mode(data, keepdims=True)[0][0]
variance = np.var(data, ddof=1)
std_dev = np.sqrt(variance)
coef_var = (std_dev / mean) * 100  # Коэффициент вариации в %

print(f"Среднее: {mean:.2f}")
print(f"Медиана: {median:.2f}")
print(f"Мода: {mode:.2f}")
print(f"Дисперсия: {variance:.2f}")
print(f"Стандартное отклонение: {std_dev:.2f}")
print(f"Коэффициент вариации: {coef_var:.2f}%")

# Оценка однородности выборки
if coef_var < 33:
    print("Выборка однородна")
else:
    print("Выборка не однородна")

# 3. Удаление выбросов (отбрасываем значения за пределами 1.5*IQR)
Q1, Q3 = np.percentile(data, [25, 75])
IQR = Q3 - Q1
lower_bound = Q1 - 1.5 * IQR
upper_bound = Q3 + 1.5 * IQR
filtered_data = data[(data >= lower_bound) & (data <= upper_bound)]

print(f"Размер выборки до удаления выбросов: {len(data)}")
print(f"Размер выборки после удаления выбросов: {len(filtered_data)}")

# 4. Разбиение на интервалы
n = len(filtered_data)
s1 = int(1 + np.log2(n))
s2 = int(1 + 3.322 * np.log10(n))

print(f"Оптимальное количество интервалов (s1): {s1}")
print(f"Оптимальное количество интервалов (s2): {s2}")

# Гистограмма и полигон
plt.figure(figsize=(12, 6))
sns.histplot(filtered_data, bins=s1, kde=True, color='blue', alpha=0.6)
plt.title("Гистограмма распределения GDP на душу населения")
plt.xlabel("GDP per capita")
plt.ylabel("Частота")
plt.show()

# 5. Дисперсионный анализ
groups = np.array_split(np.sort(filtered_data), s1)
group_means = [np.mean(g) for g in groups]
group_variances = [np.var(g, ddof=1) for g in groups]

S_within = np.mean(group_variances)  # Внутригрупповая дисперсия
S_between = np.var(group_means, ddof=1)  # Межгрупповая дисперсия
S_total = np.var(filtered_data, ddof=1)  # Общая дисперсия

print(f"Общая дисперсия: {S_total:.2f}")
print(f"Внутригрупповая дисперсия: {S_within:.2f}")
print(f"Межгрупповая дисперсия: {S_between:.2f}")

# Проверка сложения дисперсий
print(f"Сумма внутригрупповой и межгрупповой: {S_within + S_between:.2f}")
print(f"Погрешность: {abs(S_total - (S_within + S_between)):.5f}")

# 6. Проверка соответствия закону распределения
# Критерий Колмогорова-Смирнова
ks_statistic, ks_p_value = stats.kstest(filtered_data, 'norm', args=(mean, std_dev))
print(f"Критерий Колмогорова-Смирнова: Statistic={ks_statistic:.4f}, p-value={ks_p_value:.4f}")

# Критерий Пирсона (χ²)
hist, bin_edges = np.histogram(filtered_data, bins=s1)
expected_freq = len(filtered_data) / s1
chi2_stat, chi2_p_value = stats.chisquare(hist, [expected_freq] * len(hist))

print(f"Критерий Пирсона: χ²={chi2_stat:.4f}, p-value={chi2_p_value:.4f}")

# Вывод о нормальности
if ks_p_value > 0.05 and chi2_p_value > 0.05:
    print("Выборка соответствует нормальному распределению.")
else:
    print("Выборка не соответствует нормальному распределению.")