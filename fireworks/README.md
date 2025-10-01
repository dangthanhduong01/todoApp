# 🎆 Fireworks Package

Package pháo hoa với hiệu ứng vật lý thực tế cho ứng dụng Fyne.

## ✨ Tính năng

### 🎇 Hiệu ứng pháo hoa nâng cao:
- **Advanced Physics**: Particles với gravity, air resistance, wind, và bounce
- **Trail Effects**: Vệt sáng theo sau mỗi particle với fade
- **Multiple explosion types**: 5 loại nổ với mathematical precision
- **Animation mượt mà**: 5 FPS với trajectory tính toán chính xác  
- **Sophisticated Fade**: Alpha blending với fade in/out tự nhiên
- **Orchestrated Sequences**: 6 waves explosions với grand finale
- **Sparkle Effects**: 20% particles có hiệu ứng lấp lánh random
- **Bounce Physics**: 30% particles nảy khi chạm đất với energy loss

### 🎨 Các loại explosion:

1. **Burst** 🎆
   - Nổ tán ra theo mọi hướng
   - Classic fireworks explosion

2. **Fountain** ⛲ 
   - Phun lên như đài phun nước
   - Particles bay lên cao rồi rơi xuống

3. **Spiral** 🌀
   - Xoắn ốc tạo pattern đẹp mắt
   - Hiệu ứng xoay tròn

4. **Heart** 💝
   - Hình trái tim romantic
   - Sử dụng phương trình tham số toán học

5. **Star** ⭐
   - Hình ngôi sao 5 cánh
   - Pattern symmetric đẹp mắt

## 🚀 Cách sử dụng

### Import package:
```go
import "todoapp/fireworks"
```

### Hiển thị pháo hoa:
```go
fireworks.ShowFireworksDialog("Hoàn thành task!", myWindow)
```

### Tạo fireworks system custom:
```go
// Tạo system
fs := fireworks.NewFireworksSystem(800, 600)

// Thêm explosions
fs.AddExplosion(400, 300, fireworks.Burst)
fs.AddExplosion(200, 200, fireworks.Heart)

// Bắt đầu animation
fs.Start()

// Render to container
container := fs.Render()
```

## 📐 Advanced Physics Engine

### Enhanced Particle Properties:
- **Position**: X, Y coordinates với sub-pixel precision
- **Velocity**: VX, VY speed vectors với realistic physics
- **Gravity**: Downward acceleration (0.05-0.2) variable
- **Life**: Frame-based aging với MaxLife randomization
- **Size**: Variable particle sizes (0.8-2.3) với alpha scaling
- **Alpha**: Transparency (0.0-1.0) với fade in/out
- **Trail**: 5-point trail system với independent fade
- **Bounce**: 30% particles bounce với energy loss
- **Wind Resistance**: Air resistance (0.98-1.0) cho realism
- **Sparkle**: 20% particles có lấp lánh effects

### Advanced Physics Calculations:
- **Air Resistance**: `VX *= windResistance, VY *= windResistance`
- **Realistic Wind**: `windForce = 0.01 * sin(frameCount*0.1 + X*0.01)`  
- **Bounce Mechanics**: `VY *= -0.6, VX *= 0.8` khi impact
- **Turbulence**: Random drift `±0.015` cho natural movement

### Mathematical Formulas:

**Gravity Effect:**
```
VY += gravity  // Acceleration downward
Y += VY        // Update position
```

**Heart Shape (Parametric):**
```
X = 16 * sin³(t)
Y = -(13*cos(t) - 5*cos(2t) - 2*cos(3t) - cos(4t))
```

**Enhanced Star Pattern:**
```
outerRadius = explosion.Speed
innerRadius = explosion.Speed * 0.4
pointIndex = i % (starPoints * 2)
angle = pointIndex * (π / starPoints)
```

**Trail System:**
```
trail[i] = {X, Y, Alpha, Age}
trailAlpha *= 0.8  // Fade each frame  
trailSize = 10 + index*2  // Size decreases
```

**Advanced Fade Algorithm:**
```
if lifeRatio > 0.7:
    alpha = 1.0 - ((lifeRatio - 0.7) / 0.3)
elif lifeRatio < 0.1:
    alpha = lifeRatio / 0.1
```

## ⚡ HIGH-SPEED Performance

- **Frame Rate**: 10 FPS (100ms intervals) - **TĂNG GẤP ĐÔI**
- **Particle Count**: 30-70 per explosion - **NHIỀU HƠN 40%**
- **Particle Speed**: 4-8 velocity - **NHANH GẤP ĐÔI** 
- **Max Particles**: ~500 simultaneous (high density)
- **Duration**: 8 seconds rapid-fire sequence - **NHANH GẤP 3 LẦN**
- **Trail Points**: 3 points per particle (optimized cho tốc độ)
- **Gravity**: 0.15-0.4 acceleration - **MẠNH HƠN**
- **Turbulence**: Enhanced dynamic movement
- **Memory**: Ultra-fast cleanup + optimized trail management

## 🎯 Integration với Todo App

Khi user hoàn thành một task, system sẽ:

1. **Call**: `fireworks.ShowFireworksDialog(taskDescription, window)`
2. **Create**: 5 explosions với timing khác nhau
3. **Physics**: Particles bay theo trajectory thực tế
4. **Cleanup**: Auto-cleanup khi dialog đóng

## 🛠️ Technical Implementation

### Core Structures:
- `Particle`: Individual firework particle
- `Explosion`: Explosion definition  
- `FireworksSystem`: Main animation engine

### Thread Safety:
- Sử dụng `fyne.DoAndWait()` cho UI updates
- Goroutines cho physics calculations
- Proper cleanup mechanisms

## 🎨 Customization

### Colors:
```go
colors := []string{"🎆", "🎇", "✨", "🌟", "💫", "⭐", "💥", "🔥", "🌠", "🎊"}
```

### Enhanced Timing Sequence:
- **Wave 1**: Center burst (immediate)
- **Wave 2**: Symmetric bursts (+300ms)  
- **Wave 3**: Romantic heart (+900ms)
- **Wave 4**: Star constellation (+1700ms)
- **Wave 5**: Fountain finale (+2200ms)
- **Wave 6**: Spiral dance (+2600ms)
- **Grand Finale**: Massive burst (+3300ms)
- **Final Sparkles**: Random bursts (+3500ms)

## 📱 Fyne Integration

Compatible với:
- Fyne v2.6.3+
- Cross-platform (Linux, Windows, macOS)
- Software rendering support
- Responsive layouts

---

**🎊 Tận hưởng hiệu ứng pháo hoa tuyệt đẹp khi hoàn thành tasks! 🎊**
