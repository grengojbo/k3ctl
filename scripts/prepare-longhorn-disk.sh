#!/usr/bin/env bash
# Підготовка диска /dev/sdb для Longhorn на worker-ноді
# Виконується з sudo

set -euo pipefail

DISK="${LONGHORN_DISK:-/dev/sdb}"
MOUNTPOINT="/var/lib/longhorn"
LABEL="longhorn"

echo "==> Перевірка диска ${DISK}"
if [[ ! -b "${DISK}" ]]; then
  echo "ERROR: Диск ${DISK} не знайдено"
  exit 1
fi

# Перевірка чи диск вже має розділи
if lsblk "${DISK}" -n -o NAME | grep -q "${DISK##*/}[0-9]"; then
  echo "WARNING: Диск ${DISK} вже має розділи. Перевірка /var/lib/longhorn..."
  if mountpoint -q "${MOUNTPOINT}"; then
    echo "==> ${MOUNTPOINT} вже змонтовано"
    df -h "${MOUNTPOINT}"
    exit 0
  fi
fi

echo "==> Створення GPT та розділу на ${DISK}"
sudo parted -s "${DISK}" mklabel gpt
sudo parted -s "${DISK}" mkpart primary ext4 0% 100%

# Дочекатися появи розділу
sleep 1
PARTITION="${DISK}1"
if [[ ! -b "${PARTITION}" ]]; then
  echo "ERROR: Розділ ${PARTITION} не створено"
  exit 1
fi

echo "==> Створення файлової системи ext4 на ${PARTITION}"
sudo mkfs.ext4 -L "${LABEL}" "${PARTITION}"

echo "==> Створення точки монтування ${MOUNTPOINT}"
sudo mkdir -p "${MOUNTPOINT}"

echo "==> Отримання UUID"
UUID=$(sudo blkid -s UUID -o value "${PARTITION}")
if [[ -z "${UUID}" ]]; then
  echo "ERROR: Не вдалося отримати UUID для ${PARTITION}"
  exit 1
fi
echo "UUID=${UUID}"

# Перевірка чи вже є запис у fstab
if grep -q "${MOUNTPOINT}" /etc/fstab; then
  echo "WARNING: Запис для ${MOUNTPOINT} вже є у /etc/fstab"
else
  echo "==> Додавання у /etc/fstab"
  echo "UUID=${UUID} ${MOUNTPOINT} ext4 defaults,noatime,nofail 0 2" | sudo tee -a /etc/fstab
fi

echo "==> Монтування"
sudo mount -a

echo "==> Перевірка"
df -h "${MOUNTPOINT}"
lsblk "${DISK}"

echo "==> Готово: ${MOUNTPOINT} змонтовано з ${PARTITION}"
