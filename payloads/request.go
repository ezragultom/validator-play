package payloads

var request = []byte(`{
  "title": "Judul Program",
  "description": "Program bantuan untuk peningkatan kapasitas sosial masyarakat.",
  "cover_image": "https://cdn.example.com/programs/training-cover.jpg",
  "supporting_document": "https://cdn.example.com/docs/program-guide.pdf",
  "submission_deadline": "2025-12-01",
  "program_period": "Oktober - Desember 2025",
  "implementation_cost": 5000000000,
  "question": [
    {
      "question_type": "long_text",
      "label": "Nama Lengkap",
      "helper_text": "Masukkan nama sesuai KTP.",
      "attributes": {
        "placeholder": "Nama lengkap...",
        "max_char": 255,
        "input_validation": "free_text"
      },
      "showing_order": 1,
      "is_required": true
    },
    {
      "question_type": "long_text",
      "label": "Deskripsi Kegiatan",
      "helper_text": "Ceritakan kegiatan yang diajukan.",
      "attributes": {
        "placeholder": "Tulis deskripsi kegiatan Anda di sini...",
        "max_char": 500
      },
      "showing_order": 2,
      "is_required": true
    },
    {
      "question_type": "date_time",
      "label": "Tanggal Pelaksanaan",
      "helper_text": "Pilih tanggal pelaksanaan kegiatan.",
      "attributes": {
        "format": "date_time" 
      },
      "showing_order": 3,
      "is_required": true
    },
    {
      "question_type": "file_upload",
      "label": "Upload Proposal Kegiatan",
      "helper_text": "Unggah dokumen proposal dalam format PDF.",
      "attributes": {},
      "showing_order": 4,
      "is_required": true
    },
    {
      "question_type": "option",
      "label": "Jenis Bantuan yang Diajukan",
      "helper_text": "Pilih satu jenis bantuan.",
      "attributes": {
        "option_type": "radio_button",
        "choices": ["Dana Tunai", "Pelatihan", "Peralatan", "Beasiswa"]
      },
      "showing_order": 5,
      "is_required": true
    },
    {
      "question_type": "option",
      "label": "Bidang Kegiatan",
      "helper_text": "Anda boleh memilih lebih dari satu bidang.",
      "attributes": {
        "option_type": "checkbox",
        "choices": ["Pendidikan", "Kesehatan", "Lingkungan", "Ekonomi", "Sosial"]
      },
      "showing_order": 6,
      "is_required": false
    }
  ]
  
}`)
