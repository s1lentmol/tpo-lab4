import os
import csv
import matplotlib.pyplot as plt

BASE_DIR = os.path.dirname(os.path.abspath(__file__))

def get_path(filename):
    return os.path.join(BASE_DIR, filename)

def calculate_throughput_over_time(filename):
    seconds_buckets = {sec: 0 for sec in range(61)}
    
    try:
        with open(get_path(filename), mode='r') as f:
            reader = csv.DictReader(f)
            for row in reader:
                if row['status_code'] == '200':
                    sec = int(float(row['timestamp_sec']))
                    if 0 <= sec <= 60:
                        seconds_buckets[sec] += 1
    except FileNotFoundError:
        print(f"[!] Файл '{filename}' не найден.")
        return [], []

    sorted_seconds = sorted(seconds_buckets.keys())
    rps_values = [seconds_buckets[sec] for sec in sorted_seconds]
    return sorted_seconds, rps_values


def read_stress_results():
    rps = []
    p95 = []
    try:
        with open(get_path('stress_results.csv'), mode='r') as f:
            reader = csv.DictReader(f)
            for row in reader:
                rps.append(int(row['rps']))
                p95.append(int(row['p95_ms']))
    except FileNotFoundError:
        pass
    return rps, p95


def build_throughput_timeline_chart():
    fig, axs = plt.subplots(3, 1, figsize=(10, 8), sharex=True, sharey=True)
    
    configs_data = [
        (1, '#4CAF50', 'Конфигурация #1 ($2700)'), 
        (2, '#FF9800', 'Конфигурация #2 ($2900)'), 
        (3, '#F44336', 'Конфигурация #3 ($5400)')
    ]
    
    has_data = False
    
    for i, (cfg_id, color, label) in enumerate(configs_data):
        filename = f"load_cfg{cfg_id}_raw.csv"
        seconds, rps = calculate_throughput_over_time(filename)
        
        ax = axs[i]
        if seconds:
            ax.plot(seconds, rps, label=label, color=color, linewidth=1.8, alpha=0.9)
            has_data = True
            
            ax.grid(True, linestyle=':', alpha=0.6)
            ax.set_ylim(0, 11)
            ax.legend(loc='lower left', fontsize=9)
            ax.set_ylabel('RPS (Усп. запр/сек)', fontsize=9)
        else:
            ax.text(0.5, 0.5, f"{label}: нет данных", ha='center', va='center')

    if not has_data:
        plt.close()
        return

    plt.xlabel('Время от начала теста (секунды)', fontsize=11)
    fig.suptitle('Динамика пропускной способности приложения во времени (Throughput over Time)', fontsize=12, y=0.96)
    
    plt.tight_layout(rect=[0, 0, 1, 0.95])
    plt.savefig(get_path('throughput_timeline.png'), dpi=300)
    plt.close()
    print("[+] График сохранен в 'scripts/throughput_timeline.png'")


def build_stress_chart(rps, p95):
    if not rps: return
    plt.figure(figsize=(9, 5))
    
    plt.plot(rps, p95, marker='o', linestyle='-', color='#1f77b4', linewidth=2.5, label='Время отклика p95')
    plt.axhline(y=820, color='r', linestyle='--', linewidth=1.5, label='Лимит ТЗ (820 мс)')

    plt.xlabel('Нагрузка (RPS)', fontsize=11)
    plt.ylabel('Время отклика p95 (мс)', fontsize=11)
    plt.title('Зависимость времени отклика p95 от интенсивности нагрузки (Конфигурация №1)', fontsize=12, pad=15)
    plt.grid(True, linestyle=':', alpha=0.6)
    plt.legend(fontsize=10, loc='upper left')

    for i, txt in enumerate(p95):
        plt.annotate(f"{txt}ms", (rps[i], p95[i]), 
                     textcoords="offset points", 
                     xytext=(0,10), 
                     ha='center', fontsize=8)

    plt.tight_layout()
    plt.savefig(get_path('stress_test.png'), dpi=300)
    plt.close()
    print("[+] График стресс-теста сохранен в 'scripts/stress_test.png'")


if __name__ == "__main__":
    build_throughput_timeline_chart()
    
    stress_rps, stress_p95 = read_stress_results()
    build_stress_chart(stress_rps, stress_p95)