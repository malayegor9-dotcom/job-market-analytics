import streamlit as st
import pandas as pd
import plotly.express as px
import plotly.graph_objects as go
from sqlalchemy import create_engine

#  py -3.12 -m streamlit run dashboard.py

st.set_page_config(
    page_title="Job Market Analytics",
    page_icon="💼",
    layout="wide",
    initial_sidebar_state="expanded"
)

st.markdown("""
<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
* { font-family: 'Inter', sans-serif; box-sizing: border-box; }
.stApp { background: #0d0d1a; }

/* ── Сайдбар ── */
[data-testid="stSidebar"] {
    background: linear-gradient(180deg, #13102b 0%, #1a1040 100%);
    border-right: 1px solid #2a1f4e;
    min-width: 190px !important;
    max-width: 190px !important;
    width: 190px !important;
    transition: all 0.3s ease;
}
[data-testid="stSidebar"]:hover {
    box-shadow: 4px 0 24px rgba(155, 89, 245, 0.15);
}
[data-testid="stSidebar"] * { color: #c9b8ff !important; }
section[data-testid="stSidebar"] > div { padding: 14px 10px !important; }

/* Пункты меню */
[data-testid="stSidebar"] .stRadio label {
    padding: 10px 14px;
    border-radius: 12px;
    margin: 3px 0;
    display: block;
    cursor: pointer;
    transition: background 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
    border: 1px solid transparent;
}
[data-testid="stSidebar"] .stRadio label:hover {
    background: rgba(155, 89, 245, 0.15);
    transform: translateX(4px);
    border-color: #3d2b6e;
    box-shadow: 0 2px 12px rgba(155,89,245,0.15);
}

/* ── Метрики ── */
.metric-card {
    border-radius: 18px;
    padding: 20px;
    color: white;
    position: relative;
    overflow: hidden;
    transition: transform 0.25s ease, box-shadow 0.25s ease;
    cursor: default;
    margin-bottom: 4px;
}
.metric-card:hover {
    transform: translateY(-4px) scale(1.03);
    box-shadow: 0 16px 40px rgba(155, 89, 245, 0.4);
}
.metric-card::before {
    content: '';
    position: absolute;
    top: -40px; right: -40px;
    width: 120px; height: 120px;
    background: rgba(255,255,255,0.07);
    border-radius: 50%;
    transition: transform 0.3s ease;
}
.metric-card:hover::before { transform: scale(1.2); }
.metric-card::after {
    content: '';
    position: absolute;
    bottom: -25px; left: 20px;
    width: 70px; height: 70px;
    background: rgba(255,255,255,0.04);
    border-radius: 50%;
}
.metric-label { font-size: 11px; font-weight: 500; opacity: 0.75; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 8px; }
.metric-value { font-size: 32px; font-weight: 700; line-height: 1; margin-bottom: 4px; }
.metric-sub   { font-size: 11px; opacity: 0.65; }
.metric-icon  { position: absolute; top: 16px; right: 16px; font-size: 28px; opacity: 0.45; transition: transform 0.3s ease; }
.metric-card:hover .metric-icon { transform: scale(1.2) rotate(10deg); }

.card-1 { background: linear-gradient(135deg, #5c2d9e 0%, #9b59f5 100%); }
.card-2 { background: linear-gradient(135deg, #8e1a9e 0%, #d450f0 100%); }
.card-3 { background: linear-gradient(135deg, #6420c0 0%, #a86aff 100%); }
.card-4 { background: linear-gradient(135deg, #a01050 0%, #f06292 100%); }
.card-5 { background: linear-gradient(135deg, #7010a0 0%, #e040fb 100%); }

/* ── Панели ── */
.panel {
    background: linear-gradient(145deg, #13102b 0%, #1a1535 100%);
    border: 1px solid #2a1f4e;
    border-radius: 18px;
    padding: 20px;
    margin-bottom: 16px;
    transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease;
}
.panel:hover {
    transform: translateY(-3px);
    box-shadow: 0 12px 40px rgba(140, 60, 255, 0.18);
    border-color: #5c2d9e;
}
.panel-title {
    font-size: 14px;
    font-weight: 600;
    color: #c9b8ff;
    margin-bottom: 16px;
    padding-bottom: 10px;
    border-bottom: 1px solid #2a1f4e;
    display: flex;
    align-items: center;
    gap: 8px;
}
.panel-title::before {
    content: '';
    display: inline-block;
    width: 3px; height: 16px;
    background: linear-gradient(180deg, #9b59f5, #f06292);
    border-radius: 2px;
}

/* ── Вакансии ── */
.vacancy-row {
    display: flex;
    align-items: center;
    padding: 11px 14px;
    border-radius: 12px;
    gap: 12px;
    margin-bottom: 7px;
    background: linear-gradient(135deg, #1a1535 0%, #1f1845 100%);
    border: 1px solid #2a1f4e;
    transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
    cursor: pointer;
    text-decoration: none;
}
.vacancy-row:hover {
    transform: translateX(6px) scale(1.01);
    box-shadow: 0 6px 24px rgba(140, 60, 255, 0.25);
    border-color: #7b2fff;
}
.vacancy-avatar {
    width: 40px; height: 40px;
    border-radius: 12px;
    display: flex; align-items: center; justify-content: center;
    font-size: 16px; font-weight: 700; color: white; flex-shrink: 0;
    transition: transform 0.2s ease;
}
.vacancy-row:hover .vacancy-avatar { transform: scale(1.1); }
.vacancy-title   { font-weight: 600; font-size: 13px; color: #e8deff; }
.vacancy-company { font-size: 11px; color: #8b78c0; margin-top: 2px; }
.vacancy-loc     { font-size: 11px; color: #6b5a99; margin-left: auto; white-space: nowrap; }
.salary-tag      { font-size: 12px; font-weight: 600; color: #f06292; white-space: nowrap; min-width: 40px; text-align: right; }
.source-habr  { background: #2a1535; color: #f06292; padding: 2px 8px; border-radius: 8px; font-size: 10px; font-weight: 600; border: 1px solid #6b1f4e; }
.source-other { background: #1f1a35; color: #8b78c0; padding: 2px 8px; border-radius: 8px; font-size: 10px; border: 1px solid #2a1f4e; }
.status-remote { background: #2a1f4e; color: #c084fc; padding: 3px 10px; border-radius: 10px; font-size: 11px; font-weight: 600; white-space: nowrap; }
.status-office { background: #1f1a35; color: #6b5a99; padding: 3px 10px; border-radius: 10px; font-size: 11px; white-space: nowrap; }

/* ── Навыки ── */
.skill-badge {
    display: inline-block;
    padding: 5px 14px;
    border-radius: 20px;
    font-size: 12px; font-weight: 500;
    margin: 3px;
    background: linear-gradient(135deg, #1f1845, #2a1f4e);
    color: #c084fc;
    border: 1px solid #3d2b6e;
    transition: transform 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
    cursor: default;
}
.skill-badge:hover {
    transform: translateY(-3px) scale(1.08);
    box-shadow: 0 6px 20px rgba(160, 80, 255, 0.35);
    background: linear-gradient(135deg, #3d2b6e, #5c2d9e);
}

/* ── Страница ── */
.page-title { font-size: 24px; font-weight: 700; color: #e8deff; margin-bottom: 4px; }
.page-sub   { font-size: 13px; color: #6b5a99; margin-bottom: 24px; }

/* ── Инпуты ── */
.stTextInput > div > div > input {
    background: #1a1535 !important;
    border: 1px solid #2a1f4e !important;
    border-radius: 12px !important;
    color: #c9b8ff !important;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
}
.stTextInput > div > div > input:focus {
    border-color: #7b2fff !important;
    box-shadow: 0 0 0 2px rgba(123, 47, 255, 0.2) !important;
}
.stSelectbox > div > div {
    background: #1a1535 !important;
    border: 1px solid #2a1f4e !important;
    border-radius: 12px !important;
    color: #c9b8ff !important;
}

/* ── Таблица ── */
[data-testid="stDataFrame"] { border-radius: 14px; overflow: hidden; }
[data-testid="stDataFrame"] table { background: #13102b !important; }

/* ── Скрыть лишнее ── */
#MainMenu, footer, header { visibility: hidden; }
.block-container { padding-top: 1.5rem; }
[data-testid="metric-container"] { display: none; }
hr { border-color: #2a1f4e; margin: 20px 0; }

/* ── Анимации появления ── */
@keyframes fadeUp {
    from { opacity: 0; transform: translateY(20px); }
    to   { opacity: 1; transform: translateY(0); }
}
@keyframes slideRight {
    from { opacity: 0; transform: translateX(-16px); }
    to   { opacity: 1; transform: translateX(0); }
}
.metric-card { animation: fadeUp 0.4s ease both; }
.panel       { animation: fadeUp 0.5s ease both; }
.page-title  { animation: slideRight 0.4s ease both; }
</style>
""", unsafe_allow_html=True)

# ── Константы ─────────────────────────────────────────────────────────
PURPLE  = ["#6c3fc7","#9b59f5","#7b2fff","#b06aff","#8e24aa","#c084fc","#a78bfa","#7c3aed"]
PINK    = ["#c2185b","#f06292","#e91e63","#f48fb1","#ad1457","#f8bbd9","#ec4899","#db2777"]
PALETTE = PURPLE + PINK

CHART = dict(paper_bgcolor="#13102b", plot_bgcolor="#13102b",
             font_color="#c9b8ff", grid="#2a1f4e")

# ── БД ────────────────────────────────────────────────────────────────
@st.cache_resource
def get_engine():
    return create_engine(
        "postgresql+psycopg://app@127.0.0.1:5433/jobmarket"
    )

@st.cache_data(ttl=300)
def load_all():
    engine = get_engine()

    vac = pd.read_sql("""
        SELECT v.id, v.title, v.location, v.salary_min, v.salary_max,
               v.currency, v.remote, v.published_at, v.url,
               c.name as company, s.name as source
        FROM vacancies v
        LEFT JOIN companies c ON c.id = v.company_id
        LEFT JOIN sources   s ON s.id = v.source_id
        ORDER BY
            CASE WHEN s.name = 'habr' THEN 0 ELSE 1 END,
            v.published_at DESC
    """, engine)

    skills = pd.read_sql("""
        SELECT s.name, s.category, COUNT(*) as count
        FROM vacancy_skills vs
        JOIN skills s ON s.id = vs.skill_id
        GROUP BY s.name, s.category
        ORDER BY count DESC LIMIT 20
    """, engine)

    companies = pd.read_sql("""
        SELECT c.name as company, COUNT(*) as count
        FROM vacancies v
        JOIN companies c ON c.id = v.company_id
        GROUP BY c.name ORDER BY count DESC LIMIT 10
    """, engine)

    weekly = pd.read_sql("""
        SELECT DATE_TRUNC('day', published_at) as day, COUNT(*) as count
        FROM vacancies
        WHERE published_at > NOW() - INTERVAL '30 days'
        GROUP BY day ORDER BY day
    """, engine)

    sources = pd.read_sql("""
        SELECT s.name as source, COUNT(*) as count
        FROM vacancies v JOIN sources s ON s.id = v.source_id
        GROUP BY s.name
    """, engine)

    return vac, skills, companies, weekly, sources

df, skills_df, companies_df, weekly_df, sources_df = load_all()

total     = len(df)
remote    = int(df['remote'].sum())
companies = df['company'].nunique()
w_salary  = int(df['salary_min'].notna().sum())
new_today = int(weekly_df['count'].iloc[-1]) if not weekly_df.empty else 0

# ── Сайдбар ───────────────────────────────────────────────────────────
with st.sidebar:
    st.markdown("""
    <div style='text-align:center; padding:16px 0 28px;'>
        <div style='font-size:32px; animation: fadeUp 0.5s ease;'>💼</div>
        <div style='font-size:15px; font-weight:700; color:#e8deff; margin-top:8px;
             background: linear-gradient(135deg,#9b59f5,#f06292);
             -webkit-background-clip:text; -webkit-text-fill-color:transparent;'>
            Job Analytics
        </div>
        <div style='font-size:10px; color:#6b5a99; margin-top:4px; letter-spacing:1px; text-transform:uppercase;'>
            Market Dashboard
        </div>
    </div>
    """, unsafe_allow_html=True)

    page = st.radio("Навигация", [
        "📊 Обзор",
        "🔧 Навыки",
        "🏢 Компании",
        "📋 Вакансии",
    ], label_visibility="collapsed")

    st.markdown(f"""
    <div style='margin-top:28px; padding:14px; 
         background: linear-gradient(135deg, #1a1535, #2a1f4e);
         border:1px solid #3d2b6e; border-radius:14px;
         font-size:11px; color:#8b78c0; text-align:center;
         transition: transform 0.2s ease;'>
        🔄 Обновление каждые 5 мин<br><br>
        <span style='font-size:22px; font-weight:700;
              background:linear-gradient(135deg,#9b59f5,#f06292);
              -webkit-background-clip:text; -webkit-text-fill-color:transparent;'>
            {total:,}
        </span><br>
        <span style='color:#6b5a99;'>вакансий в базе</span>
    </div>
    """, unsafe_allow_html=True)

# ── Хелперы ───────────────────────────────────────────────────────────
def mchart(fig, h=220):
    fig.update_layout(
        paper_bgcolor=CHART["paper_bgcolor"],
        plot_bgcolor=CHART["plot_bgcolor"],
        font=dict(color=CHART["font_color"], size=11),
        margin=dict(l=0, r=10, t=10, b=0),
        height=h,
    )
    return fig

def mcard(col, cls, icon, label, val, sub):
    with col:
        st.markdown(f"""
        <div class="metric-card {cls}">
            <div class="metric-icon">{icon}</div>
            <div class="metric-label">{label}</div>
            <div class="metric-value">{val}</div>
            <div class="metric-sub">{sub}</div>
        </div>""", unsafe_allow_html=True)

def vacancy_card(row, i):
    avatars = PURPLE + PINK
    color  = avatars[i % len(avatars)]
    letter = (row['company'] or 'U')[0].upper()
    loc    = (row.get('location') or '—').strip().rstrip(',')[:22]
    scls   = 'status-remote' if row['remote'] else 'status-office'
    stxt   = '🌍 Remote' if row['remote'] else '🏢 Офис'
    salary = f"${int(row['salary_min']):,}" if pd.notna(row.get('salary_min')) else '—'
    source = row.get('source', '')
    src_cls = 'source-habr' if source == 'habr' else 'source-other'
    src_lbl = '🔴 Habr' if source == 'habr' else source.upper()
    url = row.get('url') or '#'

    return f"""
    <a href="{url}" target="_blank" style="text-decoration:none;">
    <div class="vacancy-row">
        <div class="vacancy-avatar" style="background:linear-gradient(135deg,{color},{color}99)">{letter}</div>
        <div style="flex:1; min-width:0;">
            <div class="vacancy-title" style="white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">{row['title']}</div>
            <div class="vacancy-company">{row['company'] or '—'}</div>
        </div>
        <div class="vacancy-loc">{loc}</div>
        <div class="salary-tag">{salary}</div>
        <span class="{scls}">{stxt}</span>
        <span class="{src_cls}">{src_lbl}</span>
    </div>
    </a>"""

# ══════════════════════════════════════════════════════════════════════
# ОБЗОР
# ══════════════════════════════════════════════════════════════════════
if "Обзор" in page:
    st.markdown('<p class="page-title">📊 Обзор рынка вакансий</p>', unsafe_allow_html=True)
    st.markdown(f'<p class="page-sub">Актуальные данные · {total:,} вакансий в базе</p>', unsafe_allow_html=True)

    c1,c2,c3,c4,c5 = st.columns(5)
    mcard(c1,"card-1","💼","Вакансий",      f"{total:,}",   "в базе данных")
    mcard(c2,"card-2","🏢","Компаний",      f"{companies:,}","уникальных")
    mcard(c3,"card-3","🌍","Удалённых",     f"{remote:,}",  f"{remote*100//max(total,1)}% от всех")
    mcard(c4,"card-4","💰","С зарплатой",   f"{w_salary:,}","указана сумма")
    mcard(c5,"card-5","🆕","Новых сегодня", f"{new_today:,}","за последний день")

    st.markdown("<br>", unsafe_allow_html=True)
    col1, col2 = st.columns([3,2])

    with col1:
        st.markdown('<div class="panel"><div class="panel-title">Динамика публикаций</div>', unsafe_allow_html=True)
        if not weekly_df.empty:
            fig = go.Figure()
            fig.add_trace(go.Scatter(
                x=weekly_df['day'], y=weekly_df['count'],
                fill='tozeroy',
                fillcolor='rgba(108,63,199,0.18)',
                line=dict(color='#9b59f5', width=2.5),
                mode='lines+markers',
                marker=dict(size=5, color='#c084fc', line=dict(color='#0d0d1a', width=1.5)),
            ))
            fig = mchart(fig)
            fig.update_layout(
                xaxis=dict(showgrid=False, tickfont=dict(size=10, color='#6b5a99')),
                yaxis=dict(gridcolor=CHART["grid"], tickfont=dict(size=10, color='#6b5a99')),
            )
            st.plotly_chart(fig, use_container_width=True)
        st.markdown('</div>', unsafe_allow_html=True)

    with col2:
        st.markdown('<div class="panel"><div class="panel-title">Remote vs Офис</div>', unsafe_allow_html=True)
        fig2 = go.Figure(go.Pie(
            labels=['Remote','Офис'], values=[remote, total-remote],
            hole=0.65,
            marker=dict(colors=['#9b59f5','#2a1f4e'], line=dict(color='#0d0d1a', width=3)),
            textinfo='none',
        ))
        pct = remote*100//max(total,1)
        fig2.add_annotation(text=f"<b>{pct}%</b>", x=0.5, y=0.55, showarrow=False,
            font=dict(size=22, color='#e8deff'))
        fig2.add_annotation(text="Remote", x=0.5, y=0.38, showarrow=False,
            font=dict(size=11, color='#8b78c0'))
        fig2 = mchart(fig2)
        fig2.update_layout(showlegend=True,
            legend=dict(orientation='h', x=0.05, y=-0.08, font=dict(size=10, color='#8b78c0')))
        st.plotly_chart(fig2, use_container_width=True)
        st.markdown('</div>', unsafe_allow_html=True)

    col3, col4 = st.columns([1,2])
    with col3:
        if not sources_df.empty:
            st.markdown('<div class="panel"><div class="panel-title">Источники</div>', unsafe_allow_html=True)
            fig3 = go.Figure(go.Pie(
                labels=sources_df['source'], values=sources_df['count'],
                hole=0.55,
                marker=dict(colors=PURPLE[:len(sources_df)], line=dict(color='#0d0d1a', width=2)),
                textinfo='label+percent', textfont=dict(size=11, color='white'),
            ))
            fig3 = mchart(fig3, h=200)
            fig3.update_layout(showlegend=False)
            st.plotly_chart(fig3, use_container_width=True)
            st.markdown('</div>', unsafe_allow_html=True)

    with col4:
        st.markdown('<div class="panel"><div class="panel-title">Последние вакансии</div>', unsafe_allow_html=True)
        html = ""
        for i, row in df.head(6).iterrows():
            html += vacancy_card(row, i)
        st.markdown(html, unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

# ══════════════════════════════════════════════════════════════════════
# НАВЫКИ
# ══════════════════════════════════════════════════════════════════════
elif "Навыки" in page:
    st.markdown('<p class="page-title">🔧 Анализ навыков</p>', unsafe_allow_html=True)
    st.markdown('<p class="page-sub">Востребованные технологии на рынке</p>', unsafe_allow_html=True)

    if skills_df.empty:
        st.markdown('<div class="panel"><p style="color:#6b5a99;">Данных по навыкам пока нет</p></div>', unsafe_allow_html=True)
    else:
        col1, col2 = st.columns(2)

        with col1:
            st.markdown('<div class="panel"><div class="panel-title">Топ навыков</div>', unsafe_allow_html=True)
            fig = px.bar(
                skills_df, x='count', y='name', orientation='h',
                color='count',
                color_continuous_scale=[[0,'#2a1f4e'],[0.4,'#7b2fff'],[1,'#f06292']],
                labels={'count':'','name':''},
            )
            fig = mchart(fig, h=420)
            fig.update_layout(
                yaxis=dict(
                    categoryorder='total ascending',
                    gridcolor='rgba(42,31,78,0.5)',
                    tickfont=dict(size=12, color='#c9b8ff'),
                    showgrid=False,
                ),
                xaxis=dict(
                    gridcolor='rgba(42,31,78,0.4)',
                    zeroline=False,
                    tickfont=dict(size=10, color='#6b5a99'),
                    title=dict(text='Вакансий', font=dict(size=11, color='#6b5a99')),
                ),
                coloraxis_showscale=False,
                bargap=0.3,
                plot_bgcolor='rgba(0,0,0,0)',
                paper_bgcolor='rgba(0,0,0,0)',
            )
            fig.update_traces(
                marker_line_width=0,
                marker_cornerradius=6,
            )
            st.plotly_chart(fig, use_container_width=True)
            st.markdown('</div>', unsafe_allow_html=True)

        with col2:
            st.markdown('<div class="panel"><div class="panel-title">Карта технологий</div>', unsafe_allow_html=True)
            cat_colors = {c: PALETTE[i % len(PALETTE)] for i, c in enumerate(skills_df['category'].unique())}
            fig2 = px.treemap(
                skills_df, path=['category','name'], values='count',
                color='category', color_discrete_map=cat_colors,
            )
            fig2.update_layout(
                paper_bgcolor='rgba(0,0,0,0)',
                margin=dict(l=0, r=0, t=0, b=0),
                height=420,
                font=dict(color='white', size=12),
            )
            fig2.update_traces(
                marker=dict(
                    line=dict(width=3, color='#0d0d1a'),
                    cornerradius=8,
                ),
                textfont=dict(color='white', size=12),
                hovertemplate='<b>%{label}</b><br>Вакансий: %{value}<extra></extra>',
            )
            st.plotly_chart(fig2, use_container_width=True)
            st.markdown('</div>', unsafe_allow_html=True)

        # Статистика по категориям
        st.markdown('<div class="panel"><div class="panel-title">По категориям</div>', unsafe_allow_html=True)
        cats = skills_df.groupby('category')['count'].sum().reset_index().sort_values('count', ascending=False)
        cat_cols = st.columns(len(cats))
        for idx, (_, cat_row) in enumerate(cats.iterrows()):
            color = PALETTE[idx % len(PALETTE)]
            with cat_cols[idx]:
                st.markdown(f"""
                <div style="
                    background: linear-gradient(135deg, {color}22, {color}44);
                    border: 1px solid {color}66;
                    border-radius: 14px;
                    padding: 16px;
                    text-align: center;
                    transition: transform 0.2s ease, box-shadow 0.2s ease;
                    cursor: default;
                " onmouseover="this.style.transform='translateY(-4px)';this.style.boxShadow='0 8px 24px {color}44'"
                   onmouseout="this.style.transform='';this.style.boxShadow=''">
                    <div style="font-size:22px; font-weight:700; color:{color};">{int(cat_row['count'])}</div>
                    <div style="font-size:11px; color:#8b78c0; margin-top:4px; text-transform:uppercase; letter-spacing:0.5px;">{cat_row['category']}</div>
                </div>
                """, unsafe_allow_html=True)
        st.markdown('</div>', unsafe_allow_html=True)

        # Бейджи навыков
        st.markdown('<div class="panel"><div class="panel-title">Все навыки</div>', unsafe_allow_html=True)
        html = "".join(
            f'<span class="skill-badge">{r["name"]} <b style="color:#f06292">{r["count"]}</b></span>'
            for _, r in skills_df.iterrows()
        )
        st.markdown(html, unsafe_allow_html=True)
        st.markdown("<br></div>", unsafe_allow_html=True)

# ══════════════════════════════════════════════════════════════════════
# КОМПАНИИ
# ══════════════════════════════════════════════════════════════════════
elif "Компании" in page:
    st.markdown('<p class="page-title">🏢 Компании</p>', unsafe_allow_html=True)
    st.markdown('<p class="page-sub">Топ работодателей по количеству вакансий</p>', unsafe_allow_html=True)

    # Метрики компаний
    c1, c2, c3 = st.columns(3)
    top_company = companies_df.iloc[0] if not companies_df.empty else None
    with c1:
        st.markdown(f"""
        <div class="metric-card card-1">
            <div class="metric-icon">🏆</div>
            <div class="metric-label">Топ компания</div>
            <div class="metric-value" style="font-size:20px;">{top_company['company'] if top_company is not None else '—'}</div>
            <div class="metric-sub">{int(top_company['count']) if top_company is not None else 0} вакансий</div>
        </div>""", unsafe_allow_html=True)
    with c2:
        st.markdown(f"""
        <div class="metric-card card-2">
            <div class="metric-icon">🏢</div>
            <div class="metric-label">Всего компаний</div>
            <div class="metric-value">{companies:,}</div>
            <div class="metric-sub">уникальных работодателей</div>
        </div>""", unsafe_allow_html=True)
    with c3:
        avg_per_company = round(total / max(companies, 1), 1)
        st.markdown(f"""
        <div class="metric-card card-4">
            <div class="metric-icon">📊</div>
            <div class="metric-label">Среднее вакансий</div>
            <div class="metric-value">{avg_per_company}</div>
            <div class="metric-sub">на компанию</div>
        </div>""", unsafe_allow_html=True)

    st.markdown("<br>", unsafe_allow_html=True)
    col1, col2 = st.columns([3, 2])

    with col1:
        st.markdown('<div class="panel"><div class="panel-title">Топ компаний</div>', unsafe_allow_html=True)
        fig = go.Figure(go.Bar(
            x=companies_df['count'],
            y=companies_df['company'],
            orientation='h',
            marker=dict(
                color=PALETTE[:len(companies_df)],
                line=dict(width=0),
                cornerradius=8,
            ),
            text=companies_df['count'],
            textposition='outside',
            textfont=dict(color='#c9b8ff', size=12),
        ))
        fig = mchart(fig, h=400)
        fig.update_layout(
            plot_bgcolor='rgba(0,0,0,0)',
            paper_bgcolor='rgba(0,0,0,0)',
            yaxis=dict(
                categoryorder='total ascending',
                gridcolor='rgba(0,0,0,0)',
                tickfont=dict(size=12, color='#c9b8ff'),
            ),
            xaxis=dict(
                showticklabels=False,
                zeroline=False,
                gridcolor='rgba(42,31,78,0.3)',
            ),
        )
        st.plotly_chart(fig, use_container_width=True)
        st.markdown('</div>', unsafe_allow_html=True)

    with col2:
        st.markdown('<div class="panel"><div class="panel-title">Доля рынка</div>', unsafe_allow_html=True)
        fig2 = go.Figure(go.Pie(
            labels=companies_df['company'],
            values=companies_df['count'],
            hole=0.6,
            marker=dict(
                colors=PALETTE[:len(companies_df)],
                line=dict(color='#0d0d1a', width=3),
            ),
            textinfo='percent',
            textfont=dict(size=11, color='white'),
            hovertemplate='<b>%{label}</b><br>%{value} вакансий<br>%{percent}<extra></extra>',
        ))
        top = companies_df.iloc[0]
        fig2.add_annotation(
            text=f"<b>{top['company'].split()[0]}</b>",
            x=0.5, y=0.55, showarrow=False,
            font=dict(size=13, color='#e8deff'),
        )
        fig2.add_annotation(
            text=f"{int(top['count'])} вак.",
            x=0.5, y=0.38, showarrow=False,
            font=dict(size=11, color='#8b78c0'),
        )
        fig2 = mchart(fig2, h=400)
        fig2.update_layout(
            paper_bgcolor='rgba(0,0,0,0)',
            legend=dict(
                font=dict(size=10, color='#8b78c0'),
                orientation='v',
                x=1.02,
                bgcolor='rgba(0,0,0,0)',
            ),
        )
        st.plotly_chart(fig2, use_container_width=True)
        st.markdown('</div>', unsafe_allow_html=True)

    # Карточки компаний
    st.markdown('<div class="panel"><div class="panel-title">Все компании из топ-10</div>', unsafe_allow_html=True)
    rows = [companies_df.iloc[i:i+5] for i in range(0, min(10, len(companies_df)), 5)]
    for row_df in rows:
        cols = st.columns(5)
        for idx, (_, comp) in enumerate(row_df.iterrows()):
            color = PALETTE[idx % len(PALETTE)]
            letter = comp['company'][0].upper()
            with cols[idx]:
                st.markdown(f"""
                <div style="
                    background: linear-gradient(135deg, #1a1535, #1f1845);
                    border: 1px solid {color}55;
                    border-radius: 16px;
                    padding: 16px 12px;
                    text-align: center;
                    margin-bottom: 8px;
                    transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
                    cursor: default;
                " onmouseover="this.style.transform='translateY(-4px)';this.style.boxShadow='0 8px 24px {color}33';this.style.borderColor='{color}'"
                   onmouseout="this.style.transform='';this.style.boxShadow='';this.style.borderColor='{color}55'">
                    <div style="
                        width:44px; height:44px; border-radius:12px;
                        background:linear-gradient(135deg,{color},{color}99);
                        display:flex; align-items:center; justify-content:center;
                        font-size:18px; font-weight:700; color:white;
                        margin: 0 auto 10px;
                    ">{letter}</div>
                    <div style="font-size:12px; font-weight:600; color:#e8deff; 
                         white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">
                        {comp['company']}
                    </div>
                    <div style="font-size:11px; color:{color}; margin-top:4px; font-weight:600;">
                        {int(comp['count'])} вакансий
                    </div>
                </div>
                """, unsafe_allow_html=True)
    st.markdown('</div>', unsafe_allow_html=True)

    
# ══════════════════════════════════════════════════════════════════════
# ВАКАНСИИ
# ══════════════════════════════════════════════════════════════════════
elif "Вакансии" in page:
    st.markdown('<p class="page-title">📋 Все вакансии</p>', unsafe_allow_html=True)
    st.markdown('<p class="page-sub">Habr Career вверху · нажми на вакансию чтобы открыть</p>', unsafe_allow_html=True)

    col1, col2, col3 = st.columns(3)
    with col1:
        search = st.text_input("🔍 Поиск", placeholder="Python Developer...")
    with col2:
        remote_f = st.selectbox("🌍 Тип", ["Все","Remote","Офис"])
    with col3:
        source_f = st.selectbox("📡 Источник", ["Все","habr","remoteok"])

    filtered = df.copy()
    if search:
        filtered = filtered[filtered['title'].str.contains(search, case=False, na=False)]
    if remote_f == "Remote":
        filtered = filtered[filtered['remote'] == True]
    elif remote_f == "Офис":
        filtered = filtered[filtered['remote'] == False]
    if source_f != "Все":
        filtered = filtered[filtered['source'] == source_f]

    st.markdown(f'<p style="color:#6b5a99;font-size:12px;margin:8px 0 12px;">Найдено: <b style="color:#c084fc">{len(filtered):,}</b> вакансий</p>', unsafe_allow_html=True)

    st.markdown('<div class="panel">', unsafe_allow_html=True)
    html = ""
    for i, (_, row) in enumerate(filtered.head(50).iterrows()):
        html += vacancy_card(row, i)
    st.markdown(html, unsafe_allow_html=True)
    st.markdown('</div>', unsafe_allow_html=True)