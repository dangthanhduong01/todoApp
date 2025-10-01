package fireworks

import (
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TrailPoint đại diện cho một điểm trong vệt sáng
type TrailPoint struct {
	X, Y  float64
	Alpha float64
	Age   int
}

// Particle đại diện cho một hạt pháo hoa
type Particle struct {
	X, Y           float64 // Vị trí hiện tại
	VX, VY         float64 // Vận tốc theo trục X và Y
	StartX, StartY float64 // Vị trí ban đầu
	Life           int     // Tuổi thọ (frames)
	MaxLife        int     // Tuổi thọ tối đa
	Size           float64 // Kích thước
	Color          string  // Màu sắc (emoji)
	Gravity        float64 // Ảnh hưởng của trọng lực
	Fade           bool    // Có fade out không
	Alpha          float64 // Độ trong suốt (0.0-1.0)
	Trail          []TrailPoint // Vệt sáng phía sau
	Bounce         bool    // Có nảy khi chạm đất không
	WindResistance float64 // Sức cản gió
	Sparkle        bool    // Có lấp lánh không
}

// FireworksSystem quản lý toàn bộ hệ thống pháo hoa
type FireworksSystem struct {
	particles  []Particle
	explosions []Explosion
	canvas     *fyne.Container
	labels     map[int]*widget.Label
	width      float64
	height     float64
	isRunning  bool
	ticker     *time.Ticker
	frameCount int
	maxFrames  int
}

// Explosion đại diện cho một vụ nổ pháo hoa
type Explosion struct {
	X, Y          float64
	ParticleCount int
	Colors        []string
	Speed         float64
	Life          int
	Type          ExplosionType
}

// ExplosionType định nghĩa các loại nổ khác nhau
type ExplosionType int

const (
	Burst    ExplosionType = iota // Nổ tán ra
	Fountain                      // Phun như đài phun nước
	Spiral                        // Xoắn ốc
	Heart                         // Hình trái tim
	Star                          // Hình ngôi sao
)

// NewFireworksSystem tạo hệ thống pháo hoa mới
func NewFireworksSystem(width, height float64) *FireworksSystem {
	return &FireworksSystem{
		particles:  make([]Particle, 0),
		explosions: make([]Explosion, 0),
		canvas:     container.NewWithoutLayout(),
		labels:     make(map[int]*widget.Label),
		width:      width,
		height:     height,
		maxFrames:  80, // 8 seconds at 10 FPS - Nhanh gấp đôi
	}
}

// AddExplosion thêm một vụ nổ mới
func (fs *FireworksSystem) AddExplosion(x, y float64, explosionType ExplosionType) {
	colors := []string{"🎆", "🎇", "✨", "🌟", "💫", "⭐", "💥", "🔥", "🌠", "🎊"}

	explosion := Explosion{
		X:             x,
		Y:             y,
		ParticleCount: 30 + rand.Intn(40), // 30-70 particles - Nhiều hơn
		Colors:        colors,
		Speed:         4.0 + rand.Float64()*4.0, // 4-8 speed - Nhanh gấp đôi
		Life:          20 + rand.Intn(15),       // 20-35 frames - Ngắn hơn để nhanh hơn
		Type:          explosionType,
	}

	fs.explosions = append(fs.explosions, explosion)
	fs.createParticlesFromExplosion(explosion)
}

// createParticlesFromExplosion tạo particles từ vụ nổ
func (fs *FireworksSystem) createParticlesFromExplosion(explosion Explosion) {
	for i := 0; i < explosion.ParticleCount; i++ {
		var vx, vy float64

		switch explosion.Type {
		case Burst:
			// Nổ tán ra theo mọi hướng
			angle := rand.Float64() * 2 * math.Pi
			speed := explosion.Speed * (0.5 + rand.Float64()*0.5)
			vx = math.Cos(angle) * speed
			vy = math.Sin(angle) * speed

		case Fountain:
			// Phun lên như đài phun nước
			angle := -math.Pi/2 + (rand.Float64()-0.5)*math.Pi/3 // -60° to -120°
			speed := explosion.Speed * (0.7 + rand.Float64()*0.6)
			vx = math.Cos(angle) * speed
			vy = math.Sin(angle) * speed

		case Spiral:
			// Xoắn ốc
			spiralAngle := float64(i) * 0.5
			radius := explosion.Speed * 0.8
			vx = math.Cos(spiralAngle) * radius
			vy = math.Sin(spiralAngle) * radius

		case Heart:
			// Hình trái tim hoàn hảo với phương trình Cardioid
			t := float64(i) / float64(explosion.ParticleCount) * 2 * math.Pi
			// Phương trình trái tim cải tiến
			scale := explosion.Speed * 0.15
			heartX := scale * 16 * math.Pow(math.Sin(t), 3)
			heartY := -scale * (13*math.Cos(t) - 5*math.Cos(2*t) - 2*math.Cos(3*t) - math.Cos(4*t))
			
			// Thêm animation theo thời gian để trái tim "đập"
			pulse := 1.0 + 0.3*math.Sin(float64(i)*0.5)
			vx = heartX * pulse
			vy = heartY * pulse

		case Star:
			// Hình ngôi sao 5 cánh hoàn hảo
			starPoints := 5
			outerRadius := explosion.Speed
			innerRadius := explosion.Speed * 0.4
			
			// Tạo 2 layers: outer points và inner points
			pointIndex := i % (starPoints * 2)
			angle := float64(pointIndex) * (math.Pi / float64(starPoints))
			
			if pointIndex%2 == 0 {
				// Outer points (điểm nhọn)
				vx = math.Cos(angle) * outerRadius
				vy = math.Sin(angle) * outerRadius
			} else {
				// Inner points (điểm lõm)
				vx = math.Cos(angle) * innerRadius
				vy = math.Sin(angle) * innerRadius
			}
		}

		particle := Particle{
			X:              explosion.X,
			Y:              explosion.Y,
			StartX:         explosion.X,
			StartY:         explosion.Y,
			VX:             vx,
			VY:             vy,
			Life:           0,
			MaxLife:        explosion.Life + rand.Intn(30), // Thêm random cho đa dạng (30-80 frames)
			Size:           0.8 + rand.Float64()*1.5,       // Size từ 0.8 đến 2.3
			Color:          explosion.Colors[rand.Intn(len(explosion.Colors))],
			Gravity:        0.15 + rand.Float64()*0.25, // Gravity từ 0.15 đến 0.4 - Nhanh hơn
			Fade:           true,
			Alpha:          1.0,                              // Bắt đầu với độ trong suốt full
			Trail:          make([]TrailPoint, 0, 3),         // Trail giảm xuống 3 điểm cho tốc độ
			Bounce:         rand.Float64() < 0.4,             // 40% chance bounce - Nhiều action hơn
			WindResistance: 0.95 + rand.Float64()*0.03,      // Air resistance 0.95-0.98 - Ít cản hơn
			Sparkle:        rand.Float64() < 0.2,             // 20% chance sparkle
		}

		fs.particles = append(fs.particles, particle)
	}
}

// Update cập nhật trạng thái của tất cả particles
func (fs *FireworksSystem) Update() {
	// Cập nhật particles
	activeParticles := make([]Particle, 0)

	for _, particle := range fs.particles {
		// Lưu vị trí cũ cho trail effect
		if len(particle.Trail) < cap(particle.Trail) {
			particle.Trail = append(particle.Trail, TrailPoint{
				X: particle.X, Y: particle.Y, Alpha: particle.Alpha, Age: 0,
			})
		} else if len(particle.Trail) > 0 {
			// Shift trail và thêm điểm mới
			copy(particle.Trail[:len(particle.Trail)-1], particle.Trail[1:])
			particle.Trail[len(particle.Trail)-1] = TrailPoint{
				X: particle.X, Y: particle.Y, Alpha: particle.Alpha, Age: 0,
			}
		}

		// Cập nhật age của trail points
		for j := range particle.Trail {
			particle.Trail[j].Age++
			particle.Trail[j].Alpha *= 0.8 // Fade trail
		}

		// Advanced Physics
		// 1. Áp dụng trọng lực
		particle.VY += particle.Gravity

		// 2. Air resistance (sức cản không khí)
		particle.VX *= particle.WindResistance
		particle.VY *= particle.WindResistance

		// 3. Hiệu ứng gió mạnh hơn
		windForce := 0.03 * math.Sin(float64(fs.frameCount)*0.2 + particle.X*0.02)
		particle.VX += windForce

		// 4. Turbulence mạnh cho dynamic movement
		particle.VX += (rand.Float64() - 0.5) * 0.04  // Tăng từ 0.015 lên 0.04
		particle.VY += (rand.Float64() - 0.5) * 0.025 // Tăng từ 0.01 lên 0.025

		// 5. Cập nhật vị trí
		particle.X += particle.VX
		particle.Y += particle.VY
		particle.Life++

		// 6. Bounce effect mạnh hơn khi chạm đất
		if particle.Bounce && particle.Y >= fs.height-20 && particle.VY > 0 {
			particle.VY *= -0.8  // Bounce mạnh hơn - ít energy loss
			particle.VX *= 0.9   // Ít friction hơn
			particle.Y = fs.height - 20
			// Thêm random burst khi bounce
			particle.VX += (rand.Float64() - 0.5) * 0.5
		}

		// 7. Fade effect nâng cao
		if particle.Fade {
			lifeRatio := float64(particle.Life) / float64(particle.MaxLife)
			if lifeRatio > 0.7 {
				// Fade out trong 30% cuối đời
				fadeProgress := (lifeRatio - 0.7) / 0.3
				particle.Alpha = 1.0 - fadeProgress
			} else {
				// Fade in trong 10% đầu đời
				if lifeRatio < 0.1 {
					particle.Alpha = lifeRatio / 0.1
				}
			}
		}

		// Kiểm tra boundary
		if particle.X < 0 || particle.X > fs.width ||
			particle.Y > fs.height || particle.Life > particle.MaxLife {
			// Particle chết, không thêm vào activeParticles
			continue
		}

		activeParticles = append(activeParticles, particle)

		// Cập nhật particle trong slice
		fs.particles[len(activeParticles)-1] = particle
	}

	fs.particles = activeParticles
	fs.frameCount++
}

// Render vẽ các particles lên canvas với effects nâng cao
func (fs *FireworksSystem) Render() *fyne.Container {
	// Clear existing labels
	fs.canvas.RemoveAll()
	fs.labels = make(map[int]*widget.Label)

	// Render trails trước
	for _, particle := range fs.particles {
		// Vẽ trail (vệt sáng phía sau)
		for j, trailPoint := range particle.Trail {
			if trailPoint.Alpha > 0.1 { // Chỉ vẽ khi còn đủ sáng
				trailLabel := widget.NewLabel("·")
				trailLabel.TextStyle = fyne.TextStyle{Bold: false}
				
				// Size giảm dần theo trail
				trailSize := float32(10 + j*2)
				trailLabel.Resize(fyne.NewSize(trailSize, trailSize))
				trailLabel.Move(fyne.NewPos(float32(trailPoint.X-float64(trailSize/2)), float32(trailPoint.Y-float64(trailSize/2))))
				
				fs.canvas.Add(trailLabel)
			}
		}
	}

	// Render particles chính
	for i, particle := range fs.particles {
		// Chọn emoji dựa trên life cycle và effects
		emoji := particle.Color
		lifeRatio := float64(particle.Life) / float64(particle.MaxLife)
		
		if particle.Sparkle && rand.Float64() < 0.3 {
			// Sparkle effect
			sparkleEmojis := []string{"✨", "💫", "🌟", "⭐"}
			emoji = sparkleEmojis[rand.Intn(len(sparkleEmojis))]
		} else if lifeRatio > 0.8 {
			// Giai đoạn cuối - fade particles
			if particle.Alpha < 0.5 {
				fadeEmojis := []string{"·", "˙", "°"}
				emoji = fadeEmojis[rand.Intn(len(fadeEmojis))]
			}
		} else if lifeRatio > 0.6 {
			// Giai đoạn giữa - transition effects  
			if rand.Float64() < 0.4 {
				transitionEmojis := []string{"✨", "💫", "⭐"}
				emoji = transitionEmojis[rand.Intn(len(transitionEmojis))]
			}
		}

		// Tạo label với size dựa trên particle size và alpha
		label := widget.NewLabel(emoji)
		label.TextStyle = fyne.TextStyle{Bold: particle.Alpha > 0.7}
		
		// Size thay đổi theo alpha và particle size
		displaySize := float32(particle.Size * 15 * particle.Alpha)
		if displaySize < 8 {
			displaySize = 8
		}
		label.Resize(fyne.NewSize(displaySize, displaySize))
		label.Move(fyne.NewPos(float32(particle.X-float64(displaySize/2)), float32(particle.Y-float64(displaySize/2))))

		fs.canvas.Add(label)
		fs.labels[i] = label
	}

	return fs.canvas
}

// Start bắt đầu animation
func (fs *FireworksSystem) Start() {
	if fs.isRunning {
		return
	}

	fs.isRunning = true
	fs.frameCount = 0
	fs.ticker = time.NewTicker(100 * time.Millisecond) // 10 FPS - Tăng gấp đôi tốc độ

	// RAPID FIRE explosions sequence - Tốc độ cực nhanh!
	go func() {
		// Wave 1: Lightning fast opening burst
		fs.AddExplosion(fs.width*0.5, fs.height*0.3, Burst)
		
		time.Sleep(150 * time.Millisecond) // Giảm từ 300ms xuống 150ms
		
		// Wave 2: Rapid symmetric bursts
		fs.AddExplosion(fs.width*0.3, fs.height*0.25, Burst)
		fs.AddExplosion(fs.width*0.7, fs.height*0.25, Burst)
		
		time.Sleep(250 * time.Millisecond) // Giảm từ 600ms xuống 250ms
		
		// Wave 3: Quick romantic heart
		fs.AddExplosion(fs.width*0.5, fs.height*0.4, Heart)
		
		time.Sleep(300 * time.Millisecond) // Giảm từ 800ms xuống 300ms
		
		// Wave 4: Rapid star constellation - đồng loạt
		fs.AddExplosion(fs.width*0.2, fs.height*0.2, Star)
		time.Sleep(50 * time.Millisecond)  // Rapid fire
		fs.AddExplosion(fs.width*0.8, fs.height*0.2, Star)
		time.Sleep(50 * time.Millisecond)  // Rapid fire
		fs.AddExplosion(fs.width*0.5, fs.height*0.15, Star)
		
		time.Sleep(200 * time.Millisecond) // Giảm từ 500ms xuống 200ms
		
		// Wave 5: Instant fountain duo
		fs.AddExplosion(fs.width*0.25, fs.height*0.7, Fountain)
		fs.AddExplosion(fs.width*0.75, fs.height*0.7, Fountain)
		
		time.Sleep(150 * time.Millisecond) // Giảm từ 400ms xuống 150ms
		
		// Wave 6: Rapid spiral dance
		fs.AddExplosion(fs.width*0.4, fs.height*0.5, Spiral)
		time.Sleep(100 * time.Millisecond) // Rapid fire
		fs.AddExplosion(fs.width*0.6, fs.height*0.5, Spiral)
		
		time.Sleep(250 * time.Millisecond) // Giảm từ 700ms xuống 250ms
		
		// MASSIVE Finale: Triple burst explosion
		fs.AddExplosion(fs.width*0.5, fs.height*0.4, Burst)
		time.Sleep(80 * time.Millisecond)  // Super rapid
		fs.AddExplosion(fs.width*0.45, fs.height*0.35, Burst)
		time.Sleep(80 * time.Millisecond)  // Super rapid  
		fs.AddExplosion(fs.width*0.55, fs.height*0.45, Burst)
		
		time.Sleep(100 * time.Millisecond) // Giảm từ 200ms xuống 100ms
		
		// RAPID Final sparkles - machine gun style
		for i := 0; i < 5; i++ { // Tăng từ 3 lên 5 explosions
			x := fs.width * (0.2 + rand.Float64()*0.6) // Rộng hơn: 20%-80%
			y := fs.height * (0.1 + rand.Float64()*0.4) // Cao hơn: 10%-50%
			fs.AddExplosion(x, y, Burst)
			time.Sleep(60 * time.Millisecond) // Giảm từ 150ms xuống 60ms - machine gun!
		}
	}()

	go func() {
		defer fs.ticker.Stop()
		defer func() { fs.isRunning = false }()

		for range fs.ticker.C {
			if fs.frameCount >= fs.maxFrames && len(fs.particles) == 0 {
				return
			}

			fyne.DoAndWait(func() {
				fs.Update()
				fs.Render()
			})
		}
	}()
}

// Stop dừng animation
func (fs *FireworksSystem) Stop() {
	fs.isRunning = false
	if fs.ticker != nil {
		fs.ticker.Stop()
	}
}

// IsRunning kiểm tra xem animation có đang chạy không
func (fs *FireworksSystem) IsRunning() bool {
	return fs.isRunning
}

// ShowFireworksDialog hiển thị dialog với hiệu ứng pháo hoa
func ShowFireworksDialog(todoDescription string, window fyne.Window) {
	// Truncate description nếu quá dài
	displayDescription := todoDescription
	if len(displayDescription) > 30 {
		displayDescription = displayDescription[:27] + "..."
	}

	// Tạo labels
	mainLabel := widget.NewLabel("🎉 Đã hoàn thành: " + displayDescription)
	mainLabel.Alignment = fyne.TextAlignCenter
	mainLabel.TextStyle = fyne.TextStyle{Bold: true}

	congratsLabel := widget.NewLabel("🎊🎉 CHÚC MỪNG! 🎉🎊")
	congratsLabel.Alignment = fyne.TextAlignCenter
	congratsLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Tạo fireworks system
	fireworksSystem := NewFireworksSystem(500, 300) // Width: 500, Height: 300
	fireworksCanvas := fireworksSystem.Render()
	fireworksCanvas.Resize(fyne.NewSize(500, 300))

	// Animation text
	animationLabel := widget.NewLabel("✨🌟 Tuyệt vời! Bạn đã hoàn thành một nhiệm vụ! 🌟✨")
	animationLabel.Alignment = fyne.TextAlignCenter
	animationLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Encouragement label
	encouragementLabel := widget.NewLabel("🚀 Tiếp tục phát huy! 🚀")
	encouragementLabel.Alignment = fyne.TextAlignCenter
	encouragementLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Layout chính
	content := container.NewVBox(
		congratsLabel,
		widget.NewSeparator(),
		container.NewPadded(fireworksCanvas), // Padded fireworks area
		widget.NewSeparator(),
		mainLabel,
		widget.NewSeparator(),
		animationLabel,
		encouragementLabel,
	)

	// Tạo dialog
	animationDialog := dialog.NewCustom("🎆🎇 HOÀN THÀNH! 🎇🎆", "Tuyệt vời!", content, window)
	animationDialog.Resize(fyne.NewSize(600, 500))

	// Bắt đầu animation khi dialog hiển thị
	fireworksSystem.Start()

	// Dọn dẹp khi dialog đóng
	animationDialog.SetOnClosed(func() {
		fireworksSystem.Stop()
	})

	// Hiển thị dialog
	animationDialog.Show()
}
